package dvaplugin

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type RELATION_TYPE = int

const (
	HAS_ONE   RELATION_TYPE = iota //加载一个子对象
	HAS_MANY                       //加载多个子对象/数组
	BELONG_TO                      //直接加载数据给到父节点
)

type RELARTION_NET map[string]*RelationLoader

type RelationLoader struct {
	input         interface{}   //输入
	result        interface{}   //输出
	childModel    interface{}   //子模型
	compareFunc   CompareFun    //比较方法
	fakey         string        //父
	sukey         string        //子
	cdb           *gorm.DB      //条件db
	relation_type RELATION_TYPE //关系
	err           error         //错误
	subModifyFunc SubModifyFunc //子数据修改

	Stash             RELARTION_NET //关系网
	OpenPrintErrStack bool          //是否开启打印错误堆栈
}

func (r *RelationLoader) GetInput() interface{} {
	return r.input
}

func (r *RelationLoader) GetResult() interface{} {
	return r.result
}
func (r *RelationLoader) Error() error {
	return r.err
}
func NewRelationLoader(input interface{}, OpenPrintErrStack bool) *RelationLoader {
	return &RelationLoader{
		input:             input,
		Stash:             RELARTION_NET{},
		OpenPrintErrStack: OpenPrintErrStack,
	}
}

func (r *RelationLoader) AddRelation(relation_type RELATION_TYPE, relation, fakey, sukey string,
	Child interface{}, compareFunc CompareFun, cdb *gorm.DB, anonps ...interface{}) *RelationLoader {
	var err error
	func(err *error) {
		defer catch(err, r.OpenPrintErrStack)
		rls := strings.Split(relation, ".")
		new_sr := &RelationLoader{
			childModel:    Child,
			compareFunc:   compareFunc,
			cdb:           cdb,
			fakey:         fakey,
			sukey:         sukey,
			relation_type: relation_type,
		}
		if len(anonps) > 0 {
			sf, ok := anonps[0].(func(p, s gjson.Result) (gjson.Result, gjson.Result))
			if ok {
				new_sr.subModifyFunc = sf
			}
		}
		r.setRelation(rls, new_sr)
	}(&err)
	if err != nil {
		r.err = err
	}
	return r
}

func (r *RelationLoader) setRelation(relations []string, sr *RelationLoader) *RelationLoader {
	if len(relations) > 1 {
		_r, ok := r.Stash[relations[0]]
		if !ok {
			panic(fmt.Sprintf("欲设置的前置关系[%s]不存在[%s]", relations[0], relations))
		}
		r = _r.setRelation(relations[1:], sr)
	} else {
		if r.Stash == nil {
			r.Stash = RELARTION_NET{}
		}
		r.Stash[relations[0]] = sr
	}
	return r
}

func (r *RelationLoader) LoadResult(db *gorm.DB) *RelationLoader {
	var err error
	func(err *error) {
		defer catch(err, r.OpenPrintErrStack)
		r.load(db)
	}(&err)
	if err != nil {
		r.err = err
	}
	return r
}

func (r *RelationLoader) load(db *gorm.DB) {

	input_v := VtoJson(r.input)
	r.result = r.input
	//取key
	var fakeys = map[string][]string{}
	var keysunq = map[string]map[string]bool{}
	for rk, rv := range r.Stash {
		for _, inv := range input_v.Array() {
			value := inv.Get(rv.fakey)
			if value.Value() == nil {
				continue
			}
			val := value.String()
			if val == "" {
				continue
			}
			if _, ok := keysunq[rk]; !ok {
				keysunq[rk] = map[string]bool{}
			}
			if _, ok := keysunq[rk][val]; !ok {
				keysunq[rk][val] = true
				fakeys[rk] = append(fakeys[rk], val)
			}
		}
	}
	//加载子项
	// var subcollect = map[string]interface{}{}
	for rk, keys := range fakeys {

		rv := r.Stash[rk]
		rv_mt_slice_t := reflect.SliceOf(reflect.TypeOf(rv.childModel))
		rv_silce := reflect.New(rv_mt_slice_t).Interface()
		subcq := db
		if rv.cdb != nil {
			subcq = rv.cdb
		} else {
			subcq = subcq.Model(rv.childModel)
		}
		rows, err := subcq.Where(rv.sukey+" in ?", keys).
			Rows()
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			element := reflect.New(reflect.TypeOf(rv.childModel)).Interface()
			db.ScanRows(rows, element)
			reflect.ValueOf(rv_silce).Elem().Set(reflect.Append(reflect.ValueOf(rv_silce).Elem(),
				reflect.ValueOf(element).Elem()))
		}

		//填入结果
		if len(rv.Stash) > 0 {
			rv.input = rv_silce
			rv.load(db)
		} else {
			rv.result = rv_silce
		}

		//生成结果
		if rv.compareFunc == nil {
			rv.compareFunc = func(p, s gjson.Result) bool {
				return p.Get(rv.fakey).String() != "" && p.Get(rv.fakey).String() == s.Get(rv.sukey).String()
			}
		}
		switch rv.relation_type {
		case HAS_ONE:
			r.result, _ = HasOneV2(r.result, rv.result, rk, rv.compareFunc, rv.subModifyFunc)
		case HAS_MANY:
			r.result, _ = HasManyV2(r.result, rv.result, rk, rv.compareFunc, rv.subModifyFunc)
		case BELONG_TO:
			r.result = BelongTo(r.result, rv.result, rk)
		default:
		}
	}

}

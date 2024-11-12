package dvaplugin

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Dataer struct {
	CF       CompareFun
	Smf      SubModifyFunc
	SubGroup gjson.Result

	Meta string //原始数据

	Keys    []string        //key
	Keysunq map[string]bool //去重
}

func NewDataer(meta string, cf CompareFun, smf SubModifyFunc, subGroup gjson.Result) *Dataer {
	return &Dataer{
		Meta:     meta,
		CF:       cf,
		Smf:      smf,
		SubGroup: subGroup,

		Keys:    []string{},
		Keysunq: map[string]bool{},
	}
}
func (d *Dataer) GetResult() interface{} {
	return gjson.Parse(d.Meta).Value()
}

/*
input gjson.Result 原始数据
// 要获取的key的深度参数
// 比如  body|bar 代表获取 input的body下的bar 的值列表
dig_key string
*/
func (d *Dataer) GetKeys(input gjson.Result, dig_key string) *Dataer {

	relations := strings.Split(dig_key, "|")
	var _relatin_first = relations[0]

	if input.IsArray() {

		if len(relations) > 1 {
			for _, iv := range input.Array() {
				d.GetKeys(iv, dig_key)
			}
		} else {
			for _, iv := range input.Array() {
				_v := iv.Get(_relatin_first).String()

				// 过滤掉空的 键值
				if _v != "" {
					if _, ok := d.Keysunq[_v]; !ok {
						d.Keysunq[_v] = true
						d.Keys = append(d.Keys, _v)
					}
				}
			}
		}

	} else {

		if len(relations) > 1 {
			d.GetKeys(input.Get(_relatin_first), strings.TrimPrefix(dig_key, fmt.Sprintf("%s|", _relatin_first)))
		} else {
			_v := input.Get(_relatin_first).String()

			// 过滤掉空的 键值
			if _v != "" {
				if _, ok := d.Keysunq[_v]; !ok {
					d.Keysunq[_v] = true
					d.Keys = append(d.Keys, _v)
				}
			}

		}
	}

	return d
}

func (s *Dataer) HasOne(input gjson.Result, this_key, relation string) *Dataer {

	relations := strings.Split(relation, "|")
	var _relatin_first = relations[0]
	var w_key = this_key

	if input.IsArray() {

		for k, iv := range input.Array() {

			if this_key == "" {
				w_key = fmt.Sprintf("%d.%s", k, _relatin_first)
			} else {
				w_key = fmt.Sprintf("%s.%d.%s", this_key, k, _relatin_first)
			}

			// fmt.Println(k)
			if len(relations) > 1 {
				meta := iv.Get(_relatin_first)
				relation = strings.TrimPrefix(relation, fmt.Sprintf("%s|", _relatin_first))

				s.HasOne(meta, w_key, relation)
			} else {
				meta := iv
				//最后一个，直接比较
				var match_v gjson.Result
				SliceFind(s.SubGroup.Array(), &match_v, func(sv gjson.Result) bool {
					return s.CF(meta, sv)
				})

				if s.Smf != nil {

					_iv, _match_v := s.Smf(iv, match_v)
					match_v = _match_v

					// 先把对应的 数组的元素父值替换掉
					_iv_key := fmt.Sprintf("%s.%d", this_key, k)
					if this_key == "" {
						_iv_key = fmt.Sprintf("%d", k)
					}
					s.Meta, _ = sjson.Set(s.Meta, _iv_key, _iv.Value())

					// s.Meta = _meta.String()
				}

				// fmt.Println(w_key)
				// fmt.Println(match_v.String())

				s.Meta = VSSetV(s.Meta, match_v.Value(), w_key)
				// fmt.Println(match_v.String())
			}

		}

	} else {
		if this_key == "" {
			w_key = relations[0]
		} else {
			w_key = fmt.Sprintf("%s.%s", this_key, _relatin_first)
		}

		iv := input
		if len(relations) > 1 {
			relation = strings.TrimPrefix(relation, fmt.Sprintf("%s|", _relatin_first))
			s.HasOne(iv, w_key, relation)
		} else {

			//最后一个，直接比较
			var match_v gjson.Result
			SliceFind(s.SubGroup.Array(), &match_v, func(sv gjson.Result) bool {
				return s.CF(iv, sv)
			})

			if s.Smf != nil {
				// _meta, _match_v := s.Smf(gjson.Parse(s.Meta), match_v)
				_iv, _match_v := s.Smf(iv, match_v)

				// 先把对应的 数组的元素父值替换掉
				if this_key == "" {
					s.Meta = _iv.String()
				} else {
					_iv_key := this_key
					s.Meta, _ = sjson.Set(s.Meta, _iv_key, _iv.Value())
				}
				// s.Meta = _meta.String()
				match_v = _match_v
			}

			// fmt.Println(w_key)
			// fmt.Println(match_v.String())
			s.Meta = VSSetV(s.Meta, match_v.Value(), w_key)
		}

	}

	return s

}

func (s *Dataer) HasMany(input gjson.Result, this_key, relation string) *Dataer {

	relations := strings.Split(relation, "|")
	var _relatin_first = relations[0]
	var w_key = this_key

	if input.IsArray() {

		for k, iv := range input.Array() {

			if this_key == "" {
				w_key = fmt.Sprintf("%d.%s", k, _relatin_first)
			} else {
				w_key = fmt.Sprintf("%s.%d.%s", this_key, k, _relatin_first)
			}

			// fmt.Println(k)
			if len(relations) > 1 {
				meta := iv.Get(_relatin_first)
				relation = strings.TrimPrefix(relation, fmt.Sprintf("%s|", _relatin_first))

				s.HasMany(meta, w_key, relation)
			} else {
				// meta := iv
				// //最后一个，直接比较
				var filter = make([]interface{}, 0)
				for _, sv := range s.SubGroup.Array() {
					if s.CF(iv, sv) {
						if s.Smf != nil {
							_iv, _sv := s.Smf(iv, sv)
							_iv_key := fmt.Sprintf("%s.%d", this_key, k)
							if this_key == "" {
								_iv_key = fmt.Sprintf("%d", k)
							}

							s.Meta, _ = sjson.Set(s.Meta, _iv_key, _iv.Value())

							filter = append(filter, _sv.Value())
						} else {
							filter = append(filter, sv.Value())
						}
					}
				}

				s.Meta, _ = sjson.Set(s.Meta, w_key, filter)
			}

		}

	} else {
		if this_key == "" {
			w_key = relations[0]
		} else {
			w_key = fmt.Sprintf("%s.%s", this_key, _relatin_first)
		}

		iv := input
		if len(relations) > 1 {
			relation = strings.TrimPrefix(relation, fmt.Sprintf("%s|", _relatin_first))
			s.HasMany(iv, w_key, relation)
		} else {

			//最后一个，直接比较
			var filter = make([]interface{}, 0)
			for _, sv := range s.SubGroup.Array() {
				if s.CF(iv, sv) {
					if s.Smf != nil {
						// _meta, _sv := s.Smf(gjson.Parse(s.Meta), sv)
						_iv, _sv := s.Smf(iv, sv)

						// 先把对应的 数组的元素父值替换掉
						if this_key == "" {
							s.Meta = _iv.String()
						} else {
							_iv_key := this_key
							s.Meta, _ = sjson.Set(s.Meta, _iv_key, _iv.Value())
						}

						// s.Meta = _meta.String()
						filter = append(filter, _sv.Value())
					} else {
						filter = append(filter, sv.Value())
					}
				}
			}

			s.Meta, _ = sjson.Set(s.Meta, w_key, filter)
		}

	}

	return s

}

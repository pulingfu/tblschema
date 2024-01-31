## 全文索引
MySQL的全文索引在创建时可以使用 WITH PARSER 语法指定解析器。解析器的目的是将文本数据分词成一系列词汇，以便更好地支持全文搜索。除了 WITH PARSER ngram，MySQL 还提供了其他一些内置的解析器。以下是一些常见的全文索引解析器：
WITH PARSER ngram：
ngram 解析器支持 n-gram 分词，其中 n 表示单词的长度。它将文本数据分成连续的 n 个字符组成的词。
WITH PARSER simple：
simple 解析器是 MySQL 默认的解析器，它简单地将文本按照空格和标点符号进行分词。
WITH PARSER porter：
porter 解析器使用 Porter Stemming 算法，它会将单词还原为它们的词干（stem）形式。例如，"running" 和 "ran" 可能都还原为 "run"。
WITH PARSER bigram：
bigram 解析器支持 bigram 分词，其中文本被分成两个连续的字符组成的词。
WITH PARSER udf_name：
你还可以使用自定义的解析器（User-Defined Parser）。你可以创建自己的解析器作为存储过程（Stored Procedure）或 User-Defined Function（UDF）。
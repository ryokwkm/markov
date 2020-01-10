package morphology

import (
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

//名詞を検索して返す
//文章から名詞以外のものを除外してスペース区切りにする
//スペース区切りごとを検索してそれをFavorite？
func GetMeishi(body string) []string {
	delWords := []string{}
	t := tokenizer.New()
	morphs := t.Tokenize(body)

	for _, m := range morphs {
		features := m.Features()
		if len(features) == 0 || features[0] != "名詞" {
			delWords = append(delWords, m.Surface)
		}
	}

	for _, del := range delWords {
		body = strings.Replace(body, del, " ", -1)
	}

	ret := []string{}
	for _, e := range strings.Split(body, " ") {
		//2文字以下のものは名詞とみなさない
		if len(e) > 4 {
			ret = append(ret, e)
		}
	}

	return ret
}

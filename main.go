package markov

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/kr/pretty"

	"github.com/sirupsen/logrus"

	"github.com/ikawaha/kagome/tokenizer"
)

//Markov マルコフ連鎖してランダムな文章を作る。
type Markov struct {
	dict       map[string][]string //基本的に辞書は使い回し
	betterLen  int                 //どのくらいの長さの文章にしたいか(その長さにするとは言ってない)
	retryCount int                 //文章の長さを揃えるのに何回試行錯誤するか
	log        *logrus.Logger
}

func init() {
	tokenizer.SysDic()
}

func NewMarkov(l []string, log *logrus.Logger) *Markov {
	//作成する文章の長さを設定
	const defaultBetterLen = 60
	const defaultRetryCount = 13
	ret := Markov{
		betterLen:  defaultBetterLen,
		retryCount: defaultRetryCount,
		log:        log,
	}
	ret.dict = ret.makeMarkovDict(l)
	return &ret
}

//作成する文章の長さを設定
func (h Markov) SetLength(i int) *Markov {
	h.betterLen = i
	return &h
}

func (h Markov) makeMarkovDict(l []string) map[string][]string {
	s := strings.Join(l, "\n")
	// 文章の構成に必要ない文字を削除し文章の終わりと判定
	s = strings.NewReplacer(" ", "", "\n", "。", "「", "。", "」", "。", "(", "。", ")", "。", "?", "。", "!", "。").Replace(s)

	t := tokenizer.New()
	morphs := t.Tokenize(s)
	now, next := "", ""
	EOSKEYWORD := ".\n。"
	//辞書を作る
	//dict[前の単語]=["次の単語","次の単語","次の単語","次の単語"]
	dict := make(map[string][]string)
	for i := range morphs {
		now = morphs[i].Surface

		if i+1 < len(morphs) {
			next = morphs[i+1].Surface
		} else {
			next = "EOS"
		}

		//該当キーワードが最終文字だったらEOSに置き換え、BOSに次のキーワードを追加。
		if strings.Index(EOSKEYWORD, now) != -1 {
			if dict["BOS"] == nil {
				dict["BOS"] = make([]string, 0)
			}
			if strings.Index(EOSKEYWORD, next) == -1 {
				dict["BOS"] = append(dict["BOS"], next)
			}
			continue
		}
		//次の文字が最終文字ならEOSに置き換え
		if strings.Index(EOSKEYWORD, next) != -1 {
			next = "EOS"
		}
		//辞書に追加
		if dict[now] == nil {
			dict[now] = make([]string, 0)
		}
		dict[now] = append(dict[now], next)
	}
	h.log.Printf("dict %# v", pretty.Formatter(dict))
	return dict
}

//MakeWord 辞書を元にマルコフ連鎖して文章を作る。
func (h Markov) MakeWord() string {
	rtn := ""
	line := ""
	for i := 0; i <= h.retryCount; i++ {
		s := int64(time.Now().Nanosecond())
		line = h.makeWord(s)
		h.log.Printf("%#v %#v  %#v", math.Abs(float64(len(rtn)-h.betterLen)), math.Abs(float64(len(line)-h.betterLen)), line)
		if math.Abs(float64(len(rtn)-h.betterLen)) > math.Abs(float64(len(line)-h.betterLen)) {
			rtn = line
		}
	}
	spew.Dump(len(rtn))
	return rtn
}

//辞書を元にマルコフ連鎖して文章を作る（内部用メソッド）。
func (h Markov) makeWord(s int64) string {
	rtn := ""
	rand.Seed(s)
	now := h.dict["BOS"][rand.Intn(len(h.dict["BOS"]))]
	for i := 0; i < 100; i++ {
		if now == "EOS" {
			break
		}
		rtn += now
		if len(h.dict[now]) != 0 {
			now = h.dict[now][rand.Intn(len(h.dict[now]))]
		}

	}
	return strings.Replace(rtn, "EOS", "", -1)
}

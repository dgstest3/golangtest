package horo

import (
	"strconv"
	"strings"
	"time"

	"encoding/xml"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

const (

// RUSSIAN
	ruDailyGeneralURL   = "http://img.ignio.com/r/export/utf/xml/daily/com.xml"
	ruDailyEroticURL    = "http://img.ignio.com/r/export/utf/xml/daily/ero.xml"
	ruDailyEmotionalURL = "http://img.ignio.com/r/export/utf/xml/daily/anti.xml"
	ruDailyBusinessURL  = "http://img.ignio.com/r/export/utf/xml/daily/bus.xml"
	ruDailyHealthURL    = "http://img.ignio.com/r/export/utf/xml/daily/hea.xml"
	ruDailyCookURL      = "http://img.ignio.com/r/export/utf/xml/daily/cook.xml"
	ruDailyLoveURL      = "http://img.ignio.com/r/export/utf/xml/daily/lov.xml"
	ruDailyMobileURL    = "http://img.ignio.com/r/export/utf/xml/daily/mob.xml"

// ENGLISH
	enDailyGeneralURL  = "http://www.stardm.com/daily-horoscopes/b1-daily-horoscopes.asp"
	enDailyLoveURL     = "http://www.stardm.com/daily-horoscopes/c1-daily-love-horoscopes.asp"
	enDailyMoneyURL    = "http://www.stardm.com/daily-horoscopes/d1-daily-money-horoscopes.asp"
	enDailyBusinessURL = "http://www.stardm.com/daily-horoscopes/e1-daily-business-horoscopes.asp"
	enDailyHealthURL   = "http://www.stardm.com/daily-horoscopes/UHF-english-health-and-fitness-horoscopes.asp"
	enDailyEroticURL   = "http://www.stardm.com/daily-horoscopes/UER-daily-erotic-horoscopes.asp"
	enDailyWorkURL     = "http://www.stardm.com/daily-horoscopes/UW-daily-work-horoscopes.asp"
	enDailyWomanURL    = "http://www.stardm.com/daily-horoscopes/UWO-daily-woman-horoscopes.asp"
	enDailyChineseURL  = "http://www.stardm.com/daily-horoscopes/F2-chinese-daily-horoscopes.asp"

// FRENCH

	frDailyGeneralURL  = "http://www.stardm.com/daily-horoscopes/T-french-daily-horoscopes.asp"
	frDailyLoveURL     = "http://www.stardm.com/daily-horoscopes/TL-french-daily-love-horoscopes.asp"
	frDailyHealthURL   = "http://www.stardm.com/daily-horoscopes/THF-french-health-and-fitness-horoscopes.asp"
	frDailyEroticURL   = "http://www.stardm.com/daily-horoscopes/TER-french-daily-erotic-horoscopes.asp"
	frDailyWorkURL     = "http://www.stardm.com/daily-horoscopes/TW-french-daily-work-horoscopes.asp"
	frDailyWomanURL    = "http://www.stardm.com/daily-horoscopes/TWO-french-daily-woman-horoscopes.asp"

// DEUTSCH

	deDailyGeneralURL  = "http://www.stardm.com/daily-horoscopes/N-german-daily-horoscopes.asp"
	deDailyLoveURL     = "http://www.stardm.com/daily-horoscopes/NL-german-daily-love-horoscopes.asp"
	deDailyHealthURL   = "http://www.stardm.com/daily-horoscopes/NHF-german-health-and-fitness-horoscopes.asp"
	deDailyWorkURL     = "http://www.stardm.com/daily-horoscopes/NW-german-daily-work-horoscopes.asp"

// SPANISH

	esDailyGeneralURL  = "http://www.stardm.com/daily-horoscopes/O-spanish-daily-horoscopes.asp"
	esDailyLoveURL     = "http://www.stardm.com/daily-horoscopes/OL-spanish-daily-love-horoscopes.asp"
	esDailyHealthURL   = "http://www.stardm.com/daily-horoscopes/OHF-spanish-health-and-fitness-horoscopes.asp"
	esDailyWorkURL     = "http://www.stardm.com/daily-horoscopes/OW-spanish-daily-work-horoscopes.asp"

// ITALIAN

	itDailyGeneralURL  = "http://www.stardm.com/daily-horoscopes/P-italian-daily-horoscopes.asp"
	itDailyLoveURL     = "http://www.stardm.com/daily-horoscopes/PL-italian-daily-love-horoscopes.asp"
	itDailyWorkURL     = "http://www.stardm.com/daily-horoscopes/PW-italian-daily-work-horoscopes.asp"

// JAPANESE

	jpDailyGeneralURL  = "http://www.stardm.com/daily-horoscopes/Q-japanese-daily-horoscopes.asp"
	jpDailyLoveURL     = "http://www.stardm.com/daily-horoscopes/QL-japanese-daily-love-horoscopes.asp"
	jpDailyWorkURL     = "http://www.stardm.com/daily-horoscopes/QW-japanese-daily-work-horoscopes.asp"

// CHINESE

	cnDailyGeneralURL = "http://www.stardm.com/daily-horoscopes/R-chinese-daily-horoscopes.asp"
	cnDailyLoveURL = "http://www.stardm.com/daily-horoscopes/RL-chinese-daily-love-horoscopes.asp"
	cnDailyWorkURL = "http://www.stardm.com/daily-horoscopes/RW-chinese-daily-work-horoscopes.asp"
	
// PORTUGUESE
	
	ptDailyGeneralURL = "http://www.stardm.com/daily-horoscopes/V-portuguese-brazilian-daily-horoscopes.asp"
	ptDailyLoveURL = "http://www.stardm.com/daily-horoscopes/VL-portuguese-daily-love-horoscopes.asp"
	ptDailyWorkURL = "http://www.stardm.com/daily-horoscopes/VW-portuguese-daily-work-horoscopes.asp" 
)

var DailyURLs = map[string]map[string]string{
	"general": {
		"en": enDailyGeneralURL,
		"ru": ruDailyGeneralURL,
		"fr": frDailyGeneralURL,
		"de": deDailyGeneralURL,
		"es": esDailyGeneralURL,
		"it": itDailyGeneralURL,
		"jp": jpDailyGeneralURL,
		"cn": cnDailyGeneralURL,
		"pt": ptDailyGeneralURL},
	"love": {
		"en": enDailyLoveURL,	
		"ru": ruDailyLoveURL,
		"fr": frDailyLoveURL,
		"de": deDailyLoveURL,
		"es": esDailyLoveURL,
		"it": itDailyLoveURL,
		"jp": jpDailyLoveURL,
		"cn": cnDailyLoveURL,
		"pt": ptDailyLoveURL},
	"money": {
		"en": enDailyMoneyURL,
		"ru": "",
		"fr": "",
		"de": "",
		"es": "",
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
	"business": {
		"en": enDailyBusinessURL,
		"ru": ruDailyBusinessURL,
		"fr": "",
		"de": "",
		"es": "",
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
	"health": {
		"en": enDailyHealthURL,
		"ru": ruDailyHealthURL,
		"fr": frDailyHealthURL,
		"de": deDailyHealthURL,
		"es": esDailyHealthURL,
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
	"erotic": {
		"en": enDailyEroticURL,
		"ru": ruDailyEroticURL,
		"fr": frDailyEroticURL,
		"de": "",
		"es": "",
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
	"work": {
		"en": enDailyWorkURL,
		"ru": "",
		"fr": frDailyWorkURL,
		"de": deDailyWorkURL,
		"es": esDailyWorkURL,
		"it": itDailyWorkURL,
		"jp": jpDailyWorkURL,
		"cn": cnDailyWorkURL,
		"pt": ptDailyWorkURL},
	"woman": {
		"en": enDailyWomanURL,
		"ru": "",
		"fr": frDailyWomanURL,
		"de": "",
		"es": "",
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
	"emotional": {
		"en": "",
		"ru": ruDailyEmotionalURL,
		"fr": "",
		"de": "",
		"es": "",
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
	"cook": {
		"en": "",
		"ru": ruDailyCookURL,
		"fr": "",
		"de": "",
		"es": "",
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
	"mobile": {
		"en": "",
		"ru": ruDailyMobileURL,
		"fr": "",
		"de": "",
		"es": "",
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
	"chinese": {
		"en": enDailyChineseURL,
		"ru": "",
		"fr": "",
		"de": "",
		"es": "",
		"it": "",
		"jp": "",
		"cn": "",
		"pt": ""},
}

type HoroDate struct {
	Today string `xml:"today,attr"`
}

type Horo struct {
	XMLName     xml.Name `xml:"horo"`
	Date        HoroDate `xml:"date"`
	Aries       string   `xml:"aries>today"`
	Taurus      string   `xml:"taurus>today"`
	Gemini      string   `xml:"gemini>today"`
	Cancer      string   `xml:"cancer>today"`
	Leo         string   `xml:"leo>today"`
	Virgo       string   `xml:"virgo>today"`
	Libra       string   `xml:"libra>today"`
	Scorpio     string   `xml:"scorpio>today"`
	Sagittarius string   `xml:"sagittarius>today"`
	Capricorn   string   `xml:"capricorn>today"`
	Aquarius    string   `xml:"aquarius>today"`
	Pisces      string   `xml:"pisces>today"`
}

type HoroCh struct {
	XMLName xml.Name `xml:"horo"`
	Date    HoroDate `xml:"date"`
	Rabbit  string   `xml:"rabbit>today"`
	Dragon  string   `xml:"dragon>today"`
	Snake   string   `xml:"snake>today"`
	Horse   string   `xml:"horse>today"`
	Sheep   string   `xml:"sheep>today"`
	Monkey  string   `xml:"monkey>today"`
	Rooster string   `xml:"rooster>today"`
	Dog     string   `xml:"dog>today"`
	Pig     string   `xml:"pig>today"`
	Rat     string   `xml:"rat>today"`
	Ox      string   `xml:"ox>today"`
	Tiger   string   `xml:"tiger>today"`
}

// App environment
type App struct {
	DB         *bolt.DB
	BucketName []byte
	Log        *logrus.Logger
	DBPath     string
}

func ParseEnHoroscopeURL(url string) (*Horo, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	zodiacs := []string{}
	horoscope := []string{}

	title := doc.Find("#content .post .entry h3").First().Text()

	arrDate := strings.Split(strings.TrimSpace(strings.Split(title, "-")[1]), " ")[1:]
	arrDate = append(arrDate, strconv.Itoa(time.Now().Year()))
	nowStrDate := strings.Join(arrDate, " ")

	horoDate, err := time.Parse("2 January 2006", nowStrDate)
	if err != nil {
		logrus.Warning(err)
		horoDate = time.Now()
	}

	doc.Find("#content .post .entry h4").Each(func(i int, s *goquery.Selection) {
		zodiacs = append(zodiacs, strings.TrimSpace(strings.ToLower(s.Text())))
	})

	doc.Find("#content .post .entry p").Each(func(i int, s *goquery.Selection) {
		horoscope = append(horoscope, strings.TrimSpace(s.Text()))
	})

	horoTuple := [][2]string{}

	e := zip(zodiacs[:12], horoscope[:12], &horoTuple)
	if e != nil {
		return nil, e
	}

	horoMap := map[string]string{}
	for _, item := range horoTuple {
		horoMap[item[0]] = item[1]
	}

	hd := HoroDate{Today: horoDate.Format("02.01.2006")}
	horo := &Horo{
		Date:        hd,
		Aries:       horoMap["aries"],
		Taurus:      horoMap["taurus"],
		Gemini:      horoMap["gemini"],
		Cancer:      horoMap["cancer"],
		Leo:         horoMap["leo"],
		Virgo:       horoMap["virgo"],
		Libra:       horoMap["libra"],
		Scorpio:     horoMap["scorpio"],
		Sagittarius: horoMap["sagittarius"],
		Capricorn:   horoMap["capricorn"],
		Aquarius:    horoMap["aquarius"],
		Pisces:      horoMap["pisces"],
	}

	return horo, nil
}

func ParseEnHoroscopeChineseURL(url string) (*HoroCh, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	zodiacs := []string{}
	horoscope := []string{}

	title := doc.Find("#content .post .entry h3").First().Text()

	arrDate := strings.Split(strings.TrimSpace(strings.Split(title, "-")[1]), " ")[1:]
	arrDate = append(arrDate, strconv.Itoa(time.Now().Year()))
	nowStrDate := strings.Join(arrDate, " ")

	horoDate, err := time.Parse("2 January 2006", nowStrDate)
	if err != nil {
		logrus.Warning(err)
		horoDate = time.Now()
	}

	doc.Find("#content .post .entry h4").Each(func(i int, s *goquery.Selection) {
		zodiacs = append(zodiacs, strings.TrimSpace(strings.ToLower(s.Text())))
	})

	doc.Find("#content .post .entry p").Each(func(i int, s *goquery.Selection) {
		horoscope = append(horoscope, strings.TrimSpace(s.Text()))
	})

	horoTuple := [][2]string{}

	e := zip(zodiacs[:12], horoscope[:12], &horoTuple)
	if e != nil {
		return nil, e
	}

	horoMap := map[string]string{}
	for _, item := range horoTuple {
		horoMap[item[0]] = item[1]
	}

	hd := HoroDate{Today: horoDate.Format("02.01.2006")}
	horo := &HoroCh{
		Date:    hd,
		Rabbit:  horoMap["rabbit"],
		Dragon:  horoMap["dragon"],
		Snake:   horoMap["snake"],
		Horse:   horoMap["horse"],
		Sheep:   horoMap["sheep"],
		Monkey:  horoMap["monkey"],
		Rooster: horoMap["rooster"],
		Dog:     horoMap["dog"],
		Pig:     horoMap["pig"],
		Rat:     horoMap["rat"],
		Ox:      horoMap["ox"],
		Tiger:   horoMap["tiger"],
	}

	return horo, nil
}

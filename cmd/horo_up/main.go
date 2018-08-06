package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"testore.me/horo"
)

var (
	log    = logrus.New()
	debug  bool
	addr   string
	dbpath string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&dbpath, "db", "horo.db", "Path to DB")
}

func main() {
	flag.Parse()

	// Init logger
	if debug {
		log.Level = logrus.DebugLevel
	}

	log.Debugln("Init DB")

	// Init DB
	db, err := bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	log.Debugln("Check DB bucket")

	bucketName := []byte("horoscopes")
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			b, err = tx.CreateBucket(bucketName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		return nil
	})

	log.Debugf("URLs: %s\n", horo.DailyURLs)

	for horoType, item := range horo.DailyURLs {

		if item["en"] != "" {
			logrus.Infof("Updating %s", item["en"])

			var d time.Time
			var rawXML []byte

			if horoType == "chinese" {
				horoscopeCh, err := horo.ParseEnHoroscopeChineseURL(item["en"])
				if err != nil {
					logrus.Errorf("Failed parse URL %s\n", item["en"])
				}

				rawXML, err = xml.Marshal(horoscopeCh)
				if err != nil {
					logrus.Errorf("Failed save to XML %s\n", item["en"])
				}
				d, _ = time.Parse("02.01.2006", horoscopeCh.Date.Today)
			} else {
				horoscope, err := horo.ParseEnHoroscopeURL(item["en"])

				if err != nil {
					logrus.Errorf("Failed parse URL %s\n", item["en"])
				}

				rawXML, err = xml.Marshal(horoscope)
				if err != nil {
					logrus.Errorf("Failed save to XML %s\n", item["en"])
				}
				d, _ = time.Parse("02.01.2006", horoscope.Date.Today)
			}

			dateStr := strings.TrimSpace(d.Format("20060102"))
			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket(bucketName)
				k := []byte(fmt.Sprintf("horo:en:%s:%s", horoType, dateStr))
				b.Put(k, rawXML)
				return nil
			})
		}

		if item["ru"] != "" {
			log.Infof("Updating %s", item["ru"])

			resp, err := http.Get(item["ru"])
			if err != nil {
				logrus.Warningln(err)
				continue
			}
			defer resp.Body.Close()
			rawXML, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Errorln(err)
				continue
			}

			horoscope := &horo.Horo{}
			err = xml.Unmarshal(rawXML, horoscope)
			if err != nil {
				logrus.Errorln(err)
				continue
			}

			d, _ := time.Parse("02.01.2006", horoscope.Date.Today)
			dateStr := d.Format("20060102")

			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket(bucketName)
				k := []byte(fmt.Sprintf("horo:ru:%s:%s", horoType, dateStr))
				b.Put(k, rawXML)
				return nil
			})

		}

	}

}

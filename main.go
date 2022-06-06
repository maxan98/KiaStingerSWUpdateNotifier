package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Post struct {
	Type    string
	Message string
	Date    string
}

var debugFirstRun bool

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go Robot.StartLifeCycle(wg)
	wg.Add(1)
	go updater(wg)
	wg.Wait()

}

func updater(wg *sync.WaitGroup) {
	defer wg.Done()
	debugFirstRun = true
	prevVersion, _ := getMostRecentUpdate()
	t := time.NewTicker(30 * time.Minute)
	aliveChecker := time.NewTicker(12 * time.Hour)
	for {
		select {
		case <-t.C:
			log.Infof("Fetching Updates..")
			recentVersion, err := getMostRecentUpdate()
			if err != nil {
				log.Error("unable to fetch")
			}
			if !reflect.DeepEqual(prevVersion, recentVersion) {
				log.Infof("Found Updates")
				Robot.SendAlert("New Software released!")
				prevVersion = recentVersion
				log.Info("Updated prev version")
			} else {
				log.Info("No updates yet")
			}
		case <-aliveChecker.C:
			log.Info("I am alive")
			Robot.SendAlert("I'm still alive")
		}
	}
}

func getMostRecentUpdate() (*Post, error) {

	p := &Post{}
	c := colly.NewCollector()
	c.TraceHTTP = true
	c.OnHTML("table.col-12.mb50.list.tableStyle0 ", func(e *colly.HTMLElement) {
		tbody := e.DOM.ChildrenFiltered("tbody").Children().First()
		tbody.Children().Not("td.notitop").Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 0:
				if !debugFirstRun {
					p.Type = strings.TrimSpace(selection.Text())
				} else {
					p.Type = time.Now().GoString()
				}
			case 1:
				p.Message = strings.TrimSpace(selection.Text())
			case 2:
				p.Date = strings.TrimSpace(selection.Text())
			case 3:

			default:
				log.Error("strange table")
			}
		})
	})
	err := c.Visit("https://update.kia.com/RU/RU/updateNoticeList")
	if err != nil {
		log.Error(err)
	}
	return p, nil
}

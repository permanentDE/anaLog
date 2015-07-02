package mongo

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"go.permanent.de/anaLog/v1/anaLog/logpoint"
	"go.permanent.de/anaLog/v1/anaLog/mode"
	"go.permanent.de/anaLog/v1/anaLog/state"
	"go.permanent.de/anaLog/v1/config"
)

type Adapter struct {
	session *mgo.Session
}

func GetAdapter() *Adapter {
	return &Adapter{session: getMgoSession()}
}

func (a *Adapter) Close() {
	a.session.Close()
}

func (a *Adapter) StorePoint(lp logpoint.LogPoint) error {
	return a.StorePoints(lp)
}

func (a *Adapter) StorePoints(lps ...logpoint.LogPoint) error {
	if len(lps) < 1 {
		return fmt.Errorf("%s", "No points given")
	}

	err := a.session.Ping()
	if err != nil {
		return err
	}
	c := a.session.DB(config.Mongo.Database).C("logpoints")
	var ifs []interface{}

	for _, lp := range lps {
		ifs = append(ifs, lp)
	}

	return c.Insert(ifs...)
}

func (a *Adapter) GetRecurring() (map[string]map[string]map[string]logpoint.LogPoint, error) {
	table := make(map[string]map[string]map[string]logpoint.LogPoint)
	err := a.session.Ping()
	if err != nil {
		return table, err
	}

	c := a.session.DB(config.Mongo.Database).C("logpoints")

	var lps []logpoint.LogPoint

	err = c.Find(bson.M{"mode": fmt.Sprint(mode.Recurring)}).Iter().All(&lps)
	if err != nil {
		return table, err
	}

	for _, lp := range lps {
		_, ok := table[lp.Task]
		if ok {
			_, ok = table[lp.Task][lp.RunId]
			if !ok {
				table[lp.Task][lp.RunId] = make(map[string]logpoint.LogPoint)
			}
		} else {
			table[lp.Task] = make(map[string]map[string]logpoint.LogPoint)
			table[lp.Task][lp.RunId] = make(map[string]logpoint.LogPoint)
		}

		table[lp.Task][lp.RunId][fmt.Sprint(lp.State)] = lp
	}

	return table, nil
}

func (a *Adapter) StoreAnalysisRCserial(lrc []byte) error {
	err := a.session.Ping()
	if err != nil {
		return err
	}

	c := a.session.DB(config.Mongo.Database).C("analysis")
	return c.Insert(BytesTransport(lrc))

	//return fmt.Errorf("%s", "not implemented")
}

func (a *Adapter) GetLatestAnalysisRCserial() ([]byte, error) {
	err := a.session.Ping()
	if err != nil {
		return nil, err
	}

	var tpo BytesTransportObject

	c := a.session.DB(config.Mongo.Database).C("analysis")
	err = c.Find(nil).Sort("-time").One(&tpo)
	return tpo.Bytes, err

	//return nil, fmt.Errorf("%s", "not implemented")
}

func (a *Adapter) GetEndByStart(begin logpoint.LogPoint) (logpoint.LogPoint, error) {
	ret := logpoint.LogPoint{}
	err := a.session.Ping()
	if err != nil {
		return ret, err
	}

	var lps []logpoint.LogPoint
	c := a.session.DB(config.Mongo.Database).C("logpoints")
	err = c.Find(bson.M{"runid": begin.RunId}).Iter().All(&lps)
	if err != nil {
		return ret, err
	}

	for _, lp := range lps {
		if lp.State == state.OK || lp.State == state.Failed || lp.State == state.CompletedWithError {
			ret = lp
			return ret, nil
		}
	}

	ret.State = state.Unknown
	return ret, nil
}

func (a *Adapter) GetLastBegin(taskname string) (logpoint.LogPoint, error) {
	var lp logpoint.LogPoint
	err := a.session.Ping()
	if err != nil {
		return lp, err
	}

	c := a.session.DB(config.Mongo.Database).C("logpoints")
	err = c.Find(bson.M{"task": taskname}).Sort("-time").One(&lp)
	return lp, err
}

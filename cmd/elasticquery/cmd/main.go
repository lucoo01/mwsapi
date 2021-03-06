package cmd

import (
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/elasticengine"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/result"
	"github.com/pkg/errors"
)

// Main represents the main interface of the elasticquery Command
func Main(a *Args) (res interface{}, err error) {

	// make a new connection
	c, err := connection.NewElasticConnection(a.ElasticPort, a.ElasticHost)
	if err != nil {
		err = errors.Wrap(err, "connection.NewElasticConnection failed")
		return
	}

	// connect
	err = connection.Connect(c)
	if err != nil {
		err = errors.Wrap(err, "connection.Connect failed")
		return
	}
	defer c.Close()

	// make a query
	q := &query.ElasticQuery{
		MathWebSearchIDs: a.IDs,
		Text:             a.Text,
	}

	// run count query (if requested)
	if a.Count {
		res, err = elasticengine.Count(c, q)
		err = errors.Wrap(err, "elasticengine.Count failed")
		return
	}

	{
		var res *result.Result
		var err error

		// run either the entire thing or the document query
		if !a.DocumentPhaseOnly {
			res, err = elasticengine.Run(c, q, a.From, a.Size)
			err = errors.Wrap(err, "elasticengine.Run failed")
		} else {
			res, err = elasticengine.RunDocument(c, q, a.From, a.Size)
			err = errors.Wrap(err, "elasticengine.RunDocument failed")
		}

		// normalize if requested
		if err == nil && a.Normalize && res != nil {
			res.Normalize()
		}

		return res, err
	}

}

//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2019 SeMI Holding B.V. (registered @ Dutch Chamber of Commerce no 75221632). All rights reserved.
//  LICENSE WEAVIATE OPEN SOURCE: https://www.semi.technology/playbook/playbook/contract-weaviate-OSS.html
//  LICENSE WEAVIATE ENTERPRISE: https://www.semi.technology/playbook/contract-weaviate-enterprise.html
//  CONCEPT: Bob van Luijt (@bobvanluijt)
//  CONTACT: hello@semi.technology
//

package aggregate

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/semi-technologies/weaviate/adapters/handlers/graphql/descriptions"
	"github.com/semi-technologies/weaviate/entities/aggregation"
	"github.com/semi-technologies/weaviate/entities/models"
)

func numericPropertyFields(class *models.Class, property *models.Property, prefix string) *graphql.Object {
	getMetaIntFields := graphql.Fields{
		"sum": &graphql.Field{
			Name:        fmt.Sprintf("%s%s%sSum", prefix, class.Class, property.Name),
			Description: descriptions.LocalAggregateSum,
			Type:        graphql.Float,
			Resolve:     makeResolveFieldAggregator("sum"),
		},
		"minimum": &graphql.Field{
			Name:        fmt.Sprintf("%s%s%sMinimum", prefix, class.Class, property.Name),
			Description: descriptions.LocalAggregateMin,
			Type:        graphql.Float,
			Resolve:     makeResolveFieldAggregator("minimum"),
		},
		"maximum": &graphql.Field{
			Name:        fmt.Sprintf("%s%s%sMaximum", prefix, class.Class, property.Name),
			Description: descriptions.LocalAggregateMax,
			Type:        graphql.Float,
			Resolve:     makeResolveFieldAggregator("maximum"),
		},
		"mean": &graphql.Field{
			Name:        fmt.Sprintf("%s%s%sMean", prefix, class.Class, property.Name),
			Description: descriptions.LocalAggregateMean,
			Type:        graphql.Float,
			Resolve:     makeResolveFieldAggregator("mean"),
		},
		"mode": &graphql.Field{
			Name:        fmt.Sprintf("%s%s%sMode", prefix, class.Class, property.Name),
			Description: descriptions.LocalAggregateMode,
			Type:        graphql.Float,
			Resolve:     makeResolveFieldAggregator("mode"),
		},
		"median": &graphql.Field{
			Name:        fmt.Sprintf("%s%s%sMedian", prefix, class.Class, property.Name),
			Description: descriptions.LocalAggregateMedian,
			Type:        graphql.Float,
			Resolve:     makeResolveFieldAggregator("median"),
		},
		"count": &graphql.Field{
			Name:        fmt.Sprintf("%s%s%sCount", prefix, class.Class, property.Name),
			Description: descriptions.LocalAggregateCount,
			Type:        graphql.Int,
			Resolve:     makeResolveFieldAggregator("count"),
		},
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:        fmt.Sprintf("%s%s%sObj", prefix, class.Class, property.Name),
		Fields:      getMetaIntFields,
		Description: descriptions.LocalAggregatePropertyObject,
	})
}

func nonNumericPropertyFields(class *models.Class,
	property *models.Property, prefix string) *graphql.Object {
	getMetaPointingFields := graphql.Fields{
		"count": &graphql.Field{
			Name:        fmt.Sprintf("%s%sCount", prefix, class.Class),
			Description: descriptions.LocalAggregateCount,
			Type:        graphql.Int,
			Resolve:     makeResolveFieldAggregator("count"),
		},
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:        fmt.Sprintf("%s%s%sObj", prefix, class.Class, property.Name),
		Fields:      getMetaPointingFields,
		Description: descriptions.LocalAggregatePropertyObject,
	})
}

func groupedByProperty(class *models.Class) *graphql.Object {
	classProperties := graphql.Fields{
		"path": &graphql.Field{
			Description: descriptions.LocalAggregateGroupedByGroupedByPath,
			Type:        graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				switch typed := p.Source.(type) {
				case aggregation.GroupedBy:
					return typed.Path, nil
				case map[string]interface{}:
					return typed["path"], nil
				default:
					return nil, fmt.Errorf("groupedBy field %s: unsupported type %T", "path", p.Source)
				}
			},
		},
		"value": &graphql.Field{
			Description: descriptions.LocalAggregateGroupedByGroupedByValue,
			Type:        graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				switch typed := p.Source.(type) {
				case aggregation.GroupedBy:
					return typed.Value, nil
				case map[string]interface{}:
					return typed["value"], nil
				default:
					return nil, fmt.Errorf("groupedBy field %s: unsupported type %T", "value", p.Source)
				}
			},
		},
	}

	classPropertiesObj := graphql.NewObject(graphql.ObjectConfig{
		Name:        fmt.Sprintf("LocalAggregate%sGroupedByObj", class.Class),
		Fields:      classProperties,
		Description: descriptions.LocalAggregateGroupedByObj,
	})

	return classPropertiesObj
}

func makeResolveFieldAggregator(aggregator string) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		switch typed := p.Source.(type) {
		case aggregation.Property:
			return typed.NumericalAggregations[aggregator], nil
		case map[string]interface{}:
			return typed[aggregator], nil
		default:
			return nil, fmt.Errorf("aggregator %s, unsupported type %T", aggregator, p.Source)
		}
	}
}

package factories

import (
	"github.com/bxcodec/faker/v3"
	"gohub/app/models/link"
)

func MakeLinks(count int) []link.Link {

	var objs []link.Link

	for i := 0; i < count; i++ {
		linkModel := link.Link{
			Name: faker.Username(),
			URL:  faker.URL(),
		}
		objs = append(objs, linkModel)
	}

	return objs
}

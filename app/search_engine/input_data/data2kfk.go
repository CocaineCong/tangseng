package input_data

import (
	"fmt"
)

const inputDataPath = "./movies_data.csv"

func DocData2Kfk() {
	docs := ReadFiles([]string{inputDataPath})
	for _, doc := range docs[1:] {
		doct, err := doc2Struct(doc)
		if err != nil {
			return
		}
		fmt.Println(doct.DocId, doct.Title, doct.Body)
	}

}

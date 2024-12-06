package excelio

import (
	"github.com/SuperH-0630/hdangan/src/assest"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"io"
)

func CreateTemplate(rt runtime.RunTime, writer io.Writer) error {
	_, err := writer.Write(assest.TemplateXlsx.Content())
	if err != nil {
		return err
	}
	return nil
}

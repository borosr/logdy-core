package modes

import (
	"io"
	"logdy/utils"
	"os"

	"github.com/nxadm/tail"
	"github.com/sirupsen/logrus"

	"logdy/models"
)

func FollowFiles(ch chan models.Message, files []string) {

	for _, file := range files {

		_, err := os.Stat(file)
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"path":  file,
				"error": err.Error(),
			}).Error("Following file changes failed")
			continue
		} else {
			utils.Logger.WithFields(logrus.Fields{
				"path": file,
			}).Info("Following file changes")
		}

		go func(file string) {
			t, err := tail.TailFile(
				file, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}})
			if err != nil {
				utils.Logger.WithFields(logrus.Fields{
					"path":  file,
					"error": err.Error(),
				}).Error("Following file changes failed")
			}

			for line := range t.Lines {
				produce(ch, line.Text, models.MessageTypeStdout, &models.MessageOrigin{File: file})
			}

		}(file)
	}

}
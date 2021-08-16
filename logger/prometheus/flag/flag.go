package flag

import (
	"github.com/iTrellis/common/logger/prometheus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// AddFlags adds the flags used by this package to the Kingpin application.
// To use the default Kingpin application, call AddFlags(kingpin.CommandLine)
func AddFlags(a *kingpin.Application, config *prometheus.Config) {
	a.Flag("log.level", "Only log messages with the given severity or above. One of: [debug, info, warn, error]").
		Default("info").StringVar(&config.Level)
	a.Flag("log.file-name", "Path to the log directory.").
		Default("./log/prometheus.log").StringVar(&config.FileName)
	a.Flag("log.move-file-type", "Move file type.[0-Never move files, 1-Move files every minute, 2-Move files every hour, 3-Move files every minute day]").
		Default("3").IntVar(&config.MoveFileType)
	a.Flag("log.max-length", "The size of one file, 0-Undivided file.").
		Default("1000000000").Int64Var(&config.MaxLength)
	a.Flag("log.max-backups", "The maximum number of saved files, 0-Save all files.").
		Default("10").IntVar(&config.MaxBackups)
}

package checks

import (
	"fmt"

	"github.com/davidpenn/sensu-plugins/sdk/sensu"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

// NewMySQLReplicationStatusCommand creates a cobra command "mysql-replication-status"
func NewMySQLReplicationStatusCommand() *cobra.Command {
	var host, port, user, password string
	var warning, critical int
	cmd := &cobra.Command{
		Use: "mysql-replication-status",
		Run: func(cmd *cobra.Command, args []string) {
			dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql", user, password, host, port)
			db, err := sqlx.Connect("mysql", dataSourceName)
			if err != nil {
				sensu.Exit(sensu.RuntimeError, err.Error())
			}
			defer db.Close()
			MySQLReplicationStatus(db, warning, critical)
		},
	}

	cmd.Flags().StringVar(&host, "host", "localhost", "database host")
	cmd.Flags().StringVar(&port, "port", "3306", "database port")
	cmd.Flags().StringVar(&user, "user", "", "database user")
	cmd.Flags().StringVar(&password, "password", "", "database password")
	cmd.Flags().IntVar(&warning, "warning", 900, "warning threshold for replication lag")
	cmd.Flags().IntVar(&critical, "critical", 1800, "critical threshold for replication lag")

	return cmd
}

// MySQLReplicationStatus checks the replication status
func MySQLReplicationStatus(db *sqlx.DB, warn, crit int) {
	status := struct {
		SlaveIOState        string `db:"Slave_IO_State"`
		SlaveIORunning      string `db:"Slave_IO_Running"`
		SlaveSQLRunning     string `db:"Slave_SQL_Running"`
		LastIOError         string `db:"Last_IO_Error"`
		LastSQLError        string `db:"Last_SQL_Error"`
		SecondsBehindMaster int    `db:"Seconds_Behind_Master"`
	}{}

	err := db.Unsafe().Get(&status, "SHOW SLAVE STATUS")
	if err != nil {
		sensu.ExitWithError(err)
	}

	if status.SlaveIORunning != "Yes" || status.SlaveSQLRunning != "Yes" {
		format := "Slave not running!\n States:\n Slave_IO_Running: %s\n Slave_SQL_Running: %s\n Last Error: %s\n"
		sensu.Exit(sensu.CriticalError, format, status.SlaveIORunning, status.SlaveSQLRunning, status.LastSQLError)
	}

	delay := status.SecondsBehindMaster
	message := fmt.Sprintf("replication delayed by %d", delay)
	switch {
	case delay >= crit:
		sensu.Exit(sensu.CriticalError, message)
	case delay >= warn:
		sensu.Exit(sensu.CriticalError, message)
	}
}

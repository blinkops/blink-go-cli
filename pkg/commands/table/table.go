package table

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blinkops/blink-go-cli/gen/cli"
	"github.com/blinkops/blink-go-cli/gen/client/table"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/spf13/cobra"
)

func CRUD(table string, usage, desc string) *cobra.Command {
	return makeOperationGroupTableCmd(table, usage, desc)
}

func makeOperationGroupTableCmd(tableName string, usage, desc string) *cobra.Command {
	operationGroupTableCmd := &cobra.Command{
		Use:   usage,
		Long:  desc,
		Short: desc,
	}
	operationTableFindCmd := MakeFind(tableName)
	operationGroupTableCmd.AddCommand(operationTableFindCmd)

	operationTableGetCmd := MakeGet(tableName)
	operationGroupTableCmd.AddCommand(operationTableGetCmd)

	operationTableCreateRecordCmd := MakeCreate(tableName)
	operationGroupTableCmd.AddCommand(operationTableCreateRecordCmd)

	operationTableDeleteRecordCmd := MakeDelete(tableName)
	operationGroupTableCmd.AddCommand(operationTableDeleteRecordCmd)

	operationTableUpdateRecordCmd := MakeUpdate(tableName)
	operationGroupTableCmd.AddCommand(operationTableUpdateRecordCmd)

	return operationGroupTableCmd
}

func MakeUpdate(tableName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: `Update a record`,
		RunE:  makeUpdateExec(tableName),
	}
	cmd.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")

	return cmd
}

func makeUpdateExec(tableName string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		tableService := newTableService(cmd)
		wsID, _ := cmd.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
		id, _ := cmd.Flags().GetString(consts.IDFlagName)
		params := table.NewTableUpdateRecordParams()
		params.Table = tableName
		params.WsID = wsID
		params.ID = id
		record, err := tableService.TableUpdateRecord(params, nil)
		if err != nil {
			return err
		}
		str, _ := json.Marshal(record)
		fmt.Printf("%s", string(str))
		return nil
	}
}

func MakeDelete(tableName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: `Delete a record`,
		RunE:  makeDeleteExec(tableName),
	}

	registerWorkspaceParamFlag(cmd)
	registerIDParamFlag(cmd)

	return cmd
}

func makeDeleteExec(tableName string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		tableService := newTableService(cmd)
		wsID, _ := cmd.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
		id, _ := cmd.Flags().GetString(consts.IDFlagName)

		params := table.NewTableDeleteRecordParams()
		params.Table = tableName
		params.WsID = wsID
		params.ID = id
		record, err := tableService.TableDeleteRecord(params, nil)
		if err != nil {
			return err
		}
		str, _ := json.Marshal(record)
		fmt.Printf("%s", string(str))
		return nil
	}
}

func MakeFind(tableName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: `Find a record`,
		RunE:  makeFindExec(tableName),
	}
	_ = registerQParamFlags(cmd)
	registerWorkspaceParamFlag(cmd)

	return cmd
}

func makeFindExec(tableName string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		tableService := newTableService(cmd)
		wsID, _ := cmd.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
		qFlag, _ := cmd.Flags().GetString(consts.QueryFlagName)

		params := table.NewTableFindRecordParams()
		params.Table = tableName
		params.WsID = wsID
		params.Q = qFlag

		record, err := tableService.TableFindRecord(params, nil)
		if err != nil {
			return err
		}
		str, _ := json.Marshal(record)
		fmt.Printf("%s", string(str))
		return nil
	}
}

func MakeGet(tableName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: `Get a record by id`,
		RunE:  makeGetExec(tableName),
	}
	registerWorkspaceParamFlag(cmd)
	registerIDParamFlag(cmd)
	return cmd
}

func makeGetExec(tableName string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		tableService := newTableService(cmd)
		wsID, _ := cmd.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
		id, _ := cmd.Flags().GetString(consts.IDFlagName)

		params := table.NewTableGetRecordParams()
		params.Table = tableName
		params.WsID = wsID
		params.ID = id

		record, err := tableService.TableGetRecord(params, nil)
		if err != nil {
			return err
		}
		str, _ := json.Marshal(record)
		fmt.Printf("%s", string(str))
		return nil
	}
}

// MakeCreate returns a cmd to handle operation tableCreateRecord
func MakeCreate(tableName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: `Creates a new record`,
		RunE:  makeCreateExec(tableName),
	}
	registerWorkspaceParamFlag(cmd)
	registerRecordParamFlag(cmd, tableName)

	return cmd
}

func registerRecordParamFlag(cmd *cobra.Command, tableName string) {
	_ = cmd.PersistentFlags().String(consts.RecordFlag, "", "json string for "+tableName)
}

func makeCreateExec(tableName string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		tableService := newTableService(nil)
		wsID, _ := cmd.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
		recordFlag, _ := cmd.Flags().GetString(consts.RecordFlag)

		tableCreateRecordParams := table.NewTableCreateRecordParams()
		tableCreateRecordParams.Table = tableName
		tableCreateRecordParams.WsID = wsID
		var m map[string]string
		if err := json.Unmarshal([]byte(recordFlag), &m); err != nil {
			return errors.New("failed to parse record json")
		}
		tableCreateRecordParams.Record = m

		record, err := tableService.TableCreateRecord(tableCreateRecordParams, nil)
		if err != nil {
			return err
		}
		fmt.Printf("%v", record)
		return nil
	}
}

func newTableService(cmd *cobra.Command) table.ClientService {
	client, _ := cli.MakeClient(cmd)
	return client.Table
}

func registerWorkspaceParamFlag(cmd *cobra.Command) *string {
	return cmd.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
}

func registerIDParamFlag(cmd *cobra.Command) *string {
	return cmd.PersistentFlags().String(consts.IDFlagName, "", "Required. record ID")
}

func registerQParamFlags(cmd *cobra.Command) error {
	qDescription := `query json, for example: { "limit": 25, "offset": 0, "filter": { "name": { "$gt": "k8s" } } }`
	var qFlagDefault string
	_ = cmd.PersistentFlags().String(consts.QueryFlagName, qFlagDefault, qDescription)
	return nil
}

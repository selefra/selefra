package table

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

func ShowTable(tableHeader []string, tableBody [][]string, tableFooter []string, setBorder bool) {
	data := tableBody
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tableHeader)
	if len(tableFooter) > 0 {
		table.SetFooter(tableFooter) // Add Footer
	}
	table.SetBorder(setBorder) // Set Border to false
	table.AppendBulk(data)     // Add Bulk Data
	table.Render()
}

func ShowRows(tableHeader []string, tableBody [][]string, tableFooter []string, setBorder bool) {
	fmtStr := ""
	l := 0
	for i := range tableHeader {
		if len(tableHeader[i]) > l {
			l = len(tableHeader[i])
		}
	}
	length := strconv.Itoa(l)
	tableF := "\t%" + length + "s"
	for i := range tableBody {
		fmtStr += fmt.Sprintf("\n***********Row %d**********\n\n", i)
		for j := range tableBody[i] {
			fmtStr += fmt.Sprintf(tableF+":\t%s\n", tableHeader[j], tableBody[i][j])
		}
	}
	fmt.Println(fmtStr)
}

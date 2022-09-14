package query

const TABLESQL = "SELECT tablename FROM pg_tables where schemaname = 'public' and tablename<>'pg_stat_statements'"
const COLUMNSQL = "select table_schema,table_name,column_name from information_schema.columns where table_schema='public' and table_name<>'pg_stat_statements'"

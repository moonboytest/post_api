# go-db-practice-gonews



config.go - 

В этом файле у нас записаны несколько переменных окружения:

PostgresString - это строка, описывающая данные для подключения к бд в таком формате:
"host=<Host name> port=<Port number> dbname=<Database name> user=<Username> password=<Password> sslmode=<you an use "prefer" or somthing else> connect_timeout=<in seconds>"

SchemaSQLPath - это строка, в которой записан абсолютный путь к файлу schema.sql
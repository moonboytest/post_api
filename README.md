# go-db-practice-gonews



config/config.go - 

В этом файле у нас записаны несколько переменных окружения:

PostgresString - это строка, описывающая данные для подключения к бд postgres в таком формате:
"host=<Host name> port=<Port number> dbname=<Database name> user=<Username> password=<Password> sslmode=<you an use "prefer" or somthing else> connect_timeout=<in seconds>"

SchemaSQLPath - это строка, в которой записан абсолютный путь к файлу schema.sql

MongoDBString - это строка, описывающая данные для подключения к бд mongo в таком формате:
mongodb+srv://<nickname>:<password>@cluster0.feej5hb.mongodb.net/?retryWrites=true&w=majority

MongoDBName - имя вашей бд в mongo
MongoCollectionName - название вашей коллекции 
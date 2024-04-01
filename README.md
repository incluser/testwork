План:
1)Поднять веб сервер на localhost:5656
2)Установить зависимости jwt и mongodb driver:
    go get go.mongodb.org/mongo-driver/mongo
    go get github.com/dgrijalva/jwt-go
    go get github.com/joho/godotenv
    go get github.com/google/uuid
    go get golang.org/x/crypto/bcrypt


3)Установить соединение с MongoDB (db: testwork, collection: users)
    mongodb+srv://onetouchx9:xvT6ON5ofAZ0Vxb6@cluster0.7ay61ae.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0


4)Разработать первый роут get_token:
    1: Вытаскиваем поле user_id, если ее нет выбрасываем ошибку
    2: Через утилиту generateToken создаем пары токенов сроком 15 минут и 7 дней
    3: Через другую утилиту generateSHA256Hash создаем хеш на основе refresh_token
    ps: bcrypt ругается на длину jwt токена (bcrypt: password length exceeds 72 bytes)
    4: В базу записываем user_id и refresh_hash а клиенту возвращаем пару токенов
    



5)Разработать второй роут refresh_token:
    1: Пользователь в теле запроса передает рефреш токен
    2: Находим ID пользователя через рефреш токен
    3: Создаем для этого пользователя новый рефреш токен и параллельно хэш этого токена
    4: Записываем новый хеш в базу а клиенту вернем новый токен





POST: http://localhost:5656/get_token

BODY: 
{
    "user_id:"test_id"
}

RESPONSE: 
{
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE5MTU4ODcsInVzZXJfaWQiOiJyYW5kb21faWQifQ.wwfN2xUCvs4Vo2opE7HNqGRIAPJ2-g19dT3GGKNMywg",
	"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI1MTk3ODcsInVzZXJfaWQiOiJyYW5kb21faWQifQ.H8plLIjZO1TJ5gaKzGoXfBtDJW7y8unT2V_iFecXkxg"
}



POST: http://localhost:5656/refresh_token 

body: {
	"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI1MjAzNTQsInVzZXJfaWQiOiJyYW5kb21faWQifQ.IEu8AfZr9TBtUYFnlialGGM_LN--Axcnt54_TTY5xFA"
}

RESPONSE:
{
	"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE5MTY1MTgsInVzZXJfaWQiOiJyYW5kb21faWQifQ.jkqBpYcokaFQelnRbgdY7Rmh5h0n-YrnnNjYSOwuq1c",
	"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI1MjA0MTgsInVzZXJfaWQiOiJyYW5kb21faWQifQ.OX6Bl5CMTsTi7eNs4s3pgE-jp4XoqKYQlbMF6n-2cW8"
}
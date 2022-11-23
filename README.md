# Program for depositing and withdrawing of user's money


1. To get your user's balance use:
```
http://localhost:4000/value?user="your user id"
```

Пример:

```
http://localhost:4000/value?user=roma
```

2. To get your user's balance in certain currency use:
```
http://localhost:4000/value?user="your user id"&currency="your currency type"
```

Пример:



3.  To deposit money on user's balance use:
```
curl -X POST http://localhost:4000/deposit?user="your user id"&amount="your user's amount to deposit"
```

Пример:

```
curl -X POST http://localhost:4000/deposit?user=roma&value=10000
```

4. To widthdraw money from user's balance use:
```
curl -X POST http://localhost:4000/withdraw?user="your user id"&amount="your user's amount to withdraw"
```

5. To transfer money from one user to another use:
```
curl -X POST http://localhost:4000/transfer?from_user="user's id from who we want to get money"&to_user="user's id to whom we want to put money"&amount="value of money to transfer"
```

Example:

```
curl -X POST http://localhost:4000/transfer\?from_user\=roma\&to_user\=vasya\&value\=100
```

6. To list all operations with user's balance use:
```
http://localhost:4000/operations?user="your user id"&sort="colon on based on which we want to sort, for example: value"&sort="asc or desc"
```

Example:

```
http://localhost:4000/operations?user=aaaa&sort=value&sort=desc
```

for pagination you can:

```
http://localhost:4000/operations?user="your user id"&page="page from which you want to start"&per_page="page to what you want to show transactions"
```


Example:

```
http://localhost:4000/operations?user=userID&page=1&per_page=5
```

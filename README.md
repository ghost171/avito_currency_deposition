# Program for depositing and withdrawing of user's money


1. To get your user's balance use:
```
http://localhost:4000/balance?user="your user id"
```

2. To get your user's balance in cwertain currency use:
```
http://localhost:4000/balance?user="your user id"&currency="your currency type"
```

3.  To deposit money on user's balance use:
```
curl -X POST http://localhost:4000/deposit?user="your user id"&amount="your user's amount to deposit"
```

4. To widthdraw money from user's balance use:
```
curl -X POST http://localhost:4000/withdraw?user="your user id"&amount="your user's amount to withdraw"
```

5. To transfer money from one user to another use:
```
curl -X POST http://localhost:4000/transfer?from_user="user's id from who we want to get money"&to_user="user's id to whom we want to put money"&amount="value of money to transfer"
```

6. To list all operations with user's balance use:
```
http://localhost:4000/operations?user="your user id"&sort="colon on based on which we want to sort, for example: value"&sort="asc or desc"
```


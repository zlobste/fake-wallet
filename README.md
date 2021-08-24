## Fake wallet

### Run the API
```
docker build -t fake-wallet:test .
docker-compose up -d
```
### Stop
```
docker-compose down -v
```

### Routes
/auth/register
```
{
    "name": "Nick",
    "surname": "Krainiuk"
    "email": "nkrayniuk@gmail.com",
    "password": "qwertyqwerty"
}
```

/auth/login
```
{
    "email": "nkrayniuk@gmail.com",
    "password": "qwertyqwerty"
}
```

/wallets
/wallets/transfer
```
{
    "from": "989ad5435f70692acde0d64cf4e73347833df3708a62f459f28a8d2abe7f0fbc",
    "to": "0c91a5be52f2e89197e10689995a70c9f11f4c1d63b9721946a6f8e1cb3cb64e",
    "amount": 90
}
```
/transactions/{limit}/{offset}
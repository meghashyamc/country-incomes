# Country Incomes

## Intro

* A CLI app to project incomes/amounts from one country to another 

* This app can help answer questions like these:

> If someone earns 100,000 USD living in the USA, how much would that income be equivalent to in India, taking into account purchasing power parity?

To find that out, you can run:

```
incomes project --from united\ states --to india --amount 100000 
```

> What is the average income of a person in Germany as expressed in Chinese yuan or Indian rupees, taking into account purchasing power parity in China/India?

To find that out, you can run:

```
incomes averageincome --from germany --to india
```

This will tell you what the average income in Germany would look like to an Indian in rupees (adjusted for purchasing power parity).

Or, you could run:

```
incomes averageincome --from germany --to china
```

This will tell you what the average income in Germany would look like to a Chinese person in yuan (adjusted for purchasing power parity).

## Run

As of now, this CLI has been tested only on Linux based systems.

You can run it by building and running the 'incomes' binary with options. The environment variables mentioned in .env are compulsory (except for Redis and Postgres related environment variables which are optional)

Build the binary: 

```
cd <download path>/country-incomes
go build -o "incomes"
```
Add it to your 'path'

```
sudo cp incomes /usr/bin/
```

Checkout available options:
```
incomes --help
```

## Redis and Postgres

Redis and Postgres are used for caching country GDP values per capita (PPP) and purchasing power parity conversion factors between countries. They are optional.

However, if Redis and Postgres environment variables are specified, this can speed up the performance of the CLI significantly. For example, on one Linux system, projecting the average income of Germany to India was speeded up from around 5s to 10ms after using Redis and Postgres.

Only using Redis and not using Postgres (for partial caching) is also possible.

### Postgres migrations

If Postgres is to be used, then migrations should be run first.

```
incomes --migrate
```

 


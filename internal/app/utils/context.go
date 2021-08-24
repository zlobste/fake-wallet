package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zlobste/fake-wallet/internal/data/postgres"
	"net/http"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	usersCtxKey
	transactionsCtxKey
	walletsCtxKey
	assetsCtxKey
	JWTCtxKey
)

func CtxLog(entry *logrus.Logger) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logrus.Logger {
	return r.Context().Value(logCtxKey).(*logrus.Logger)
}

func CtxJWT(auth Auth) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, JWTCtxKey, auth)
	}
}

func JWT(r *http.Request) Auth {
	return r.Context().Value(JWTCtxKey).(Auth)
}

func CtxUsers(users postgres.UsersStorage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, usersCtxKey, users)
	}
}

func Users(r *http.Request) postgres.UsersStorage {
	return r.Context().Value(usersCtxKey).(postgres.UsersStorage).New()
}

func CtxWallets(wallets postgres.WalletsStorage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, walletsCtxKey, wallets)
	}
}

func Wallets(r *http.Request) postgres.WalletsStorage {
	return r.Context().Value(walletsCtxKey).(postgres.WalletsStorage).New()
}

func CtxTransactions(wallets postgres.TransactionsStorage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, transactionsCtxKey, wallets)
	}
}

func Transactions(r *http.Request) postgres.TransactionsStorage {
	return r.Context().Value(transactionsCtxKey).(postgres.TransactionsStorage).New()
}

func CtxAssets(wallets postgres.AssetsStorage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, assetsCtxKey, wallets)
	}
}

func Assets(r *http.Request) postgres.AssetsStorage {
	return r.Context().Value(assetsCtxKey).(postgres.AssetsStorage).New()
}
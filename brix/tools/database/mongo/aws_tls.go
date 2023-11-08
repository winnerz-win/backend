package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type awsInfo struct {
	clusterEndpoint string
	username        string
	pwd             string
	caFilePath      string
	readPreference  string
	retryWrites     bool
}

func (my awsInfo) GetURLList() []string    { return []string{} }
func (my awsInfo) IsAuth() bool            { return true }
func (my awsInfo) UserName() string        { return my.username }
func (my awsInfo) PWD() string             { return my.pwd }
func (my awsInfo) AuthSource() string      { return "" }
func (my awsInfo) ClusterEndpoint() string { return my.clusterEndpoint }
func (my awsInfo) CaFilePath() string      { return my.caFilePath }
func (my awsInfo) ReadPreference() string  { return my.readPreference }

const (
	connectTimeout           = 5
	connectionStringTemplate = "mongodb://%s:%s@%s/sample-database?tls=true&replicaSet=rs0&readpreference=%s&retryWrites=%v"
)

func NewAWSTLS(clusterEndpoint, username, password, caFilePath, readPreference string, retryWrites bool) *CDB {
	fmt.Println("clusterEndpoint :", clusterEndpoint)
	fmt.Println("username :", username)
	fmt.Println("password :", password)
	fmt.Println("caFilePath :", caFilePath)
	fmt.Println("readPreference :", readPreference)

	info := awsInfo{
		clusterEndpoint: clusterEndpoint,
		username:        username,
		pwd:             password,
		caFilePath:      caFilePath,
		readPreference:  readPreference,
		retryWrites:     retryWrites,
	}

	session := newSessionAWS(info)
	if session == nil {
		return nil
	}

	return &CDB{
		dbSession:  session,
		authSource: info.AuthSource(),
	}
}
func newSessionAWS(info awsInfo) *dbSession {
	connectionURI := fmt.Sprintf(
		connectionStringTemplate,

		info.username,
		info.pwd,
		info.clusterEndpoint,
		info.readPreference,
		info.retryWrites,
	)

	tlsConfig, err := getCustomTLSConfig(info.caFilePath)
	if err != nil {
		log.Fatalf("Failed getting TLS configuration: %v", err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI).SetTLSConfig(tlsConfig))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to DocumentDB!")

	return &dbSession{
		iInfo:  info,
		client: client,
	}
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}

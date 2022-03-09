package adxchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/marsofsnow/eos-go"
	"net/http"
	"strings"
)

const retry int = 1

type AdxChainApi struct {
	EosApi       *eos.API
	eosTxOptions *eos.TxOptions
	ServerAddress string
}

func NewAdxChainApi(chainID,chainUrl,privateKey string) *AdxChainApi {

	chainId, err := hex.DecodeString(chainID)
	if err != nil {
		panic(err)
	}

	EosApi := eos.New(chainUrl)
	EosApi.EnableKeepAlives()


	EosApi.Debug = false
	EosApi.DefaultMaxCPUUsageMS = uint8(255)
	/*
		Eospubkey, err := ecc.NewPublicKey(c.EOS.Publickey)
			if err != nil {
				logx.Errorf("ecc.NewPublicKey fail,reason:%s", err.Error())
				panic(err)
				//return nil
			}
		EosApi.SetCustomGetRequiredKeys(func(ctx context.Context, tx *eos.Transaction) ([]ecc.PublicKey, error) {
			return []ecc.PublicKey{Eospubkey}, nil
		})
	*/

	txOpts := &eos.TxOptions{
		ChainID:       chainId,
		MaxCPUUsageMS: 250,
	}

	//EosApi.DefaultMaxNetUsageWords=1024
	keyBag := &eos.KeyBag{}
	err = keyBag.ImportPrivateKey(context.Background(), privateKey)
	if err != nil {
		//return shared.Status(http.StatusInternalServerError,err.Error())
		panic(fmt.Errorf("import private key: %w", err))
	}
	EosApi.SetSigner(keyBag)

	return &AdxChainApi{
		EosApi:       EosApi,
		eosTxOptions: txOpts,
		ServerAddress: chainUrl,
	}

}

func (api *AdxChainApi) PushTransaction(actions []*eos.Action) error {

	var actionNames []string
	for _,v:=range  actions{
		actionNames=append(actionNames,v.Name.String())
	}
	ctx := context.Background()


	txOpts := &eos.TxOptions{
		//ChainID: api.eosTxOptions.ChainID,
		//DelaySecs: 500,
		//MaxCPUUsageMS: 255,
	}
	if err := txOpts.FillFromChain(ctx, api.EosApi); err != nil {
		return Status(http.StatusInternalServerError, err.Error())
	}
	tx := eos.NewTransaction(actions, txOpts)
	// signedTx, packedTx, err := api.EosApi.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
	_, packedTx, err := api.EosApi.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)

	//tx := eos.NewTransaction(actions, s.EosTxOptions)
	//signedTx, packedTx, err := s.EosApi.SignTransaction(ctx, tx, s.EosTxOptions.ChainID, eos.CompressionNone)
	if err != nil {
		return Status(http.StatusInternalServerError, err.Error())
		//panic(fmt.Errorf("sign transaction: %w", err))
	}



	_,err = api.EosApi.PushTransaction(ctx, packedTx)
	if err != nil {
		return Status(http.StatusInternalServerError, err.Error())
		//panic(fmt.Errorf("push transaction: %w", err))
	}

    /*
	logx.Infof("EOS Transaction(%s)  submitted to the network succesfully.[%s]\n",
		strings.Join(actionNames,"|"), response.StatusCode)

     */
	return nil

}


func (api *AdxChainApi) getRows(req *eos.GetTableRowsRequest) (out *eos.GetTableRowsResp, err error){
	out, err = api.EosApi.GetTableRows(context.Background(), *req)
	if err != nil  {
		if strings.Contains(err.Error(),"http: server closed idle connection"){
			count:=retry
			for count >0  {
				count--
				out, err = api.EosApi.GetTableRows(context.Background(), *req)
				if err != nil && strings.Contains(err.Error(),"http: server closed idle connection"){
					continue
				}else {
					break
				}

			}

		}

	}
	return out,err
}

// v  is array
func (api *AdxChainApi) FetchOneByPk(code, table, pk string, v interface{}) error {


	req := &eos.GetTableRowsRequest{
		Code:  code,
		Scope: code,
		Table: table,
		LowerBound: pk,
		UpperBound: pk,
		JSON:       true,
	}
	out, err := api.getRows(req)
	if err != nil  {
		if strings.Contains(err.Error(),"http: server closed idle connection"){
		}
		return Status(http.StatusInternalServerError, err.Error())
	}
	err = out.JSONToStructs(v)
	if err != nil {
		return Status(http.StatusInternalServerError, err.Error())
	}

	return nil

}

func (api *AdxChainApi) GetAccountActivePublicKey(account string) (string, error) {
	out, err := api.EosApi.GetAccount(context.Background(), eos.AccountName(account))
	if err != nil {
		return "", Status(http.StatusInternalServerError, err.Error())
	}
	//Unmarshal: public key should start with ["PUB_K1_" | "PUB_R1_" | "PUB_WA_"] (or the old "EOS")

	for index := range out.Permissions {
		if out.Permissions[index].PermName == "active" {
			return out.Permissions[index].RequiredAuth.Keys[0].PublicKey.String(), nil //默认取第一个作为public key
		}

	}
	return "",Status(http.StatusForbidden, "not found user publickey")

}



func (api *AdxChainApi) FetchBatch(code, scope, table string, limit uint32, reverse bool, v interface{}) error {


	req := &eos.GetTableRowsRequest{
		Code:    code,
		Scope:   scope,
		Table:   table,
		Limit:   limit,
		Reverse: reverse,
		JSON:    true,
	}
	out, err := api.getRows(req)
	if err != nil {
		return Status(http.StatusInternalServerError, err.Error())
	}
	err = out.JSONToStructs(v)
	if err != nil {
		return Status(http.StatusInternalServerError, err.Error())
	}
	return nil

}

func (api *AdxChainApi) FetchBatchByIndex(code, scope, table string,
	lowerBound, upperBound,index,KeyType string,
	limit uint32, reverse bool, v interface{}) error {


	req := &eos.GetTableRowsRequest{
		Code:    code,
		Scope:   scope,
		Table:   table,
		LowerBound: lowerBound,
		UpperBound: upperBound,
		KeyType: KeyType,
		Index: index,
		Limit:   limit,
		Reverse: reverse,
		JSON:    true,
	}
	out, err := api.getRows(req)
	if err != nil {

		return Status(http.StatusInternalServerError, err.Error())
	}
	err = out.JSONToStructs(v)
	if err != nil {
		return Status(http.StatusInternalServerError, err.Error())
	}


	return nil

}


func (api *AdxChainApi) FetchMoreByIndex(code, scope, table string,
	lowerBound, upperBound,index,keyType string,
	limit uint32, reverse bool, v interface{}) (bool,error) {


	req := &eos.GetTableRowsRequest{
		Code:    code,
		Scope:   scope,
		Table:   table,
		LowerBound: lowerBound,
		UpperBound: upperBound,
		Index: index,
		KeyType: keyType,
		Limit:   limit,
		Reverse: reverse,
		JSON:    true,
	}
	out, err := api.getRows(req)
	if err != nil {
		return false,Status(http.StatusInternalServerError, err.Error())
	}
	err = out.JSONToStructs(v)
	if err != nil {
		return false,Status(http.StatusInternalServerError, err.Error())
	}


	return out.More,nil

}

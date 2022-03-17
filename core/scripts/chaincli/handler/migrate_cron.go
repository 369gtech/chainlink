package handler

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/cron_upkeep_factory_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/upkeep_registration_requests_wrapper"
)

func (k *Keeper) MigrateCron(ctx context.Context, inputFile string) {
	cronFactoryAddr := common.HexToAddress(k.cfg.CronFactoryAddr)
	registrarAddr := common.HexToAddress(k.cfg.RegistrarAddr)
	// TODO
	// Add permission proxy config and transaction
	// Move target function, targetHandler, cron to file
	// MOVE LINK to fund to config
	// Write output to a file
	// Implement functionality to get upkeep ID
	// Output good tx hashes
	// write output to a file migrate_cron_output.csv
	// possibly verify contract
	targetContractAddr := common.HexToAddress("0x32f8F3021F36558f7822A283cddD3C7C1Eae9071")
	const definition = `[{"type":"function","name":"handler1","inputs":[],"outputs":[]}]`
	cronAbi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		log.Fatalln(err)
	}

	targetHandler, err := cronAbi.Pack("handler1")
	if err != nil {
		log.Fatalln(err)
	}

	cronFactoryInstance, err := cron_upkeep_factory_wrapper.NewCronUpkeepFactory(
		cronFactoryAddr,
		k.client,
	)
	if err != nil {
		log.Fatal(cronFactoryAddr.Hex(), ", is not cron factory", err)
	}
	callOpts := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}
	encodedJob, err := cronFactoryInstance.EncodeCronJob(&callOpts, targetContractAddr, targetHandler, "* * * * *")
	if err != nil {
		log.Fatalln("encoding job", err)
	}

	cronJobAddrTx, err := cronFactoryInstance.NewCronUpkeepWithJob(k.buildTxOpts(ctx), encodedJob)
	if err != nil {
		log.Fatalln("creating cron job", err)
	}
	k.waitTx(ctx, cronJobAddrTx)

	txReceipt, err := k.client.TransactionReceipt(ctx, cronJobAddrTx.Hash())
	if err != nil {
		log.Fatalln("getting receipt", err)
	}

	log1 := *txReceipt.Logs[1]
	abiLog, err := cronFactoryInstance.ParseLog(log1)
	if err != nil {
		log.Fatalln("Parsing log", err)
	}
	parsedLog, ok := abiLog.(*cron_upkeep_factory_wrapper.CronUpkeepFactoryNewCronUpkeepCreated)
	if !ok {
		log.Fatalln("Parsing log")
	}

	newCronAddr := parsedLog.Upkeep

	registrarABI, err := abi.JSON(strings.NewReader(upkeep_registration_requests_wrapper.UpkeepRegistrationRequestsABI))
	if err != nil {
		log.Fatalln("Parsing abi", err)
	}
	amount, ok := new(big.Int).SetString("15000000000000000000", 10)
	if !ok {
		log.Fatalln("Parsing nbigint", err)
	}
	registrationData, err := registrarABI.Pack("register", "name", []byte{}, newCronAddr, uint32(3000000), k.fromAddr, []byte{}, amount, uint8(0))
	if err != nil {
		log.Fatalln("generating reg data", err)
	}
	registrationTx, err := k.linkToken.TransferAndCall(k.buildTxOpts(ctx), registrarAddr, amount, registrationData)
	if err != nil {
		log.Fatalln("registering", err)
	}
	k.waitTx(ctx, registrationTx)

	fmt.Println(registrationTx.Hash())
}

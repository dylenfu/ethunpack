package unpack

import (
	"errors"
	"fmt"
	"github.com/dylenfu/extractor/abi"
	"github.com/dylenfu/extractor/types"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

const (
	erc20AbiStr         = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
	wethAbiStr          = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"guy\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"}]"
	implAbiStr          = "[{\"constant\":false,\"inputs\":[{\"name\":\"addressList\",\"type\":\"address[3][]\"},{\"name\":\"uintArgsList\",\"type\":\"uint256[7][]\"},{\"name\":\"uint8ArgsList\",\"type\":\"uint8[1][]\"},{\"name\":\"buyNoMoreThanAmountBList\",\"type\":\"bool[]\"},{\"name\":\"vList\",\"type\":\"uint8[]\"},{\"name\":\"rList\",\"type\":\"bytes32[]\"},{\"name\":\"sList\",\"type\":\"bytes32[]\"},{\"name\":\"minerId\",\"type\":\"uint256\"},{\"name\":\"feeSelections\",\"type\":\"uint16\"}],\"name\":\"submitRing\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ENTERED_MASK\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nameRegistryAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MARGIN_SPLIT_PERCENTAGE_BASE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ringIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"RATE_RATIO_SCALE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lrcTokenAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelledOrFilled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenRegistryAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delegateAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token1\",\"type\":\"address\"},{\"name\":\"token2\",\"type\":\"address\"},{\"name\":\"cutoff\",\"type\":\"uint256\"}],\"name\":\"cancelAllOrdersByTradingPair\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addresses\",\"type\":\"address[4]\"},{\"name\":\"orderValues\",\"type\":\"uint256[7]\"},{\"name\":\"buyNoMoreThanAmountB\",\"type\":\"bool\"},{\"name\":\"marginSplitPercentage\",\"type\":\"uint8\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_RING_SIZE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"bytes20\"}],\"name\":\"tradingPairCutoffs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"cutoff\",\"type\":\"uint256\"}],\"name\":\"cancelAllOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token1\",\"type\":\"address\"},{\"name\":\"token2\",\"type\":\"address\"}],\"name\":\"getTradingPairCutoffs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"cutoffs\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rateRatioCVSThreshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"walletSplitPercentage\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_lrcTokenAddress\",\"type\":\"address\"},{\"name\":\"_tokenRegistryAddress\",\"type\":\"address\"},{\"name\":\"_delegateAddress\",\"type\":\"address\"},{\"name\":\"_nameRegistryAddress\",\"type\":\"address\"},{\"name\":\"_rateRatioCVSThreshold\",\"type\":\"uint256\"},{\"name\":\"_walletSplitPercentage\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_ringIndex\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"_ringHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_miner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_feeRecipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_orderHashList\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"name\":\"_amountsList\",\"type\":\"uint256[6][]\"}],\"name\":\"RingMined\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_amountCancelled\",\"type\":\"uint256\"}],\"name\":\"OrderCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_cutoff\",\"type\":\"uint256\"}],\"name\":\"AllOrdersCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_token1\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_token2\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_cutoff\",\"type\":\"uint256\"}],\"name\":\"OrdersCancelled\",\"type\":\"event\"}]"
	delegateAbiStr      = "[{\"constant\":true,\"inputs\":[],\"name\":\"latestAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"getLatestAuthorizedAddresses\",\"outputs\":[{\"name\":\"addresses\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"authorizeAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"lrcTokenAddress\",\"type\":\"address\"},{\"name\":\"feeRecipient\",\"type\":\"address\"},{\"name\":\"batch\",\"type\":\"bytes32[]\"}],\"name\":\"batchTransferToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAddressAuthorized\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"deauthorizeAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"number\",\"type\":\"uint32\"}],\"name\":\"AddressAuthorized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"number\",\"type\":\"uint32\"}],\"name\":\"AddressDeauthorized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"
	tokenRegistryAbiStr = "[{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"unregisterToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getAddressBySymbol\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addressList\",\"type\":\"address[]\"}],\"name\":\"areAllTokensRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isTokenRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TOKEN_STANDARD_ERC223\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getTokenStandard\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"getTokens\",\"outputs\":[{\"name\":\"addressList\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TOKEN_STANDARD_ERC20\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"symbol\",\"type\":\"string\"},{\"name\":\"standard\",\"type\":\"uint8\"}],\"name\":\"registerStandardToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"registerToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"addresses\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"isTokenRegisteredBySymbol\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"TokenRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"TokenUnregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"
	nameRegistryAbiStr  = "[{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"}],\"name\":\"unregisterName\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"}],\"name\":\"registerName\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"getParticipantIds\",\"outputs\":[{\"name\":\"idList\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getParticipantById\",\"outputs\":[{\"name\":\"feeRecipient\",\"type\":\"address\"},{\"name\":\"signer\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"participantId\",\"type\":\"uint256\"}],\"name\":\"removeParticipant\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"nameInfoMap\",\"outputs\":[{\"name\":\"name\",\"type\":\"bytes12\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes12\"}],\"name\":\"ownerMap\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"nameMap\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"feeRecipient\",\"type\":\"address\"},{\"name\":\"singer\",\"type\":\"address\"}],\"name\":\"addParticipant\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"participantMap\",\"outputs\":[{\"name\":\"feeRecipient\",\"type\":\"address\"},{\"name\":\"signer\",\"type\":\"address\"},{\"name\":\"name\",\"type\":\"bytes12\"},{\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"name\",\"type\":\"string\"},{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"NameRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"name\",\"type\":\"string\"},{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"NameUnregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"name\",\"type\":\"bytes12\"},{\"indexed\":false,\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransfered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"name\",\"type\":\"bytes12\"},{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"participantId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"singer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"ParticipantRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"participantId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ParticipantUnregistered\",\"type\":\"event\"}]"
)

var (
	Erc20Abi         = newAbi(erc20AbiStr)
	WethAbi          = newAbi(wethAbiStr)
	ImplAbi          = newAbi(implAbiStr)
	DelegateAbi      = newAbi(delegateAbiStr)
	TokenRegistryAbi = newAbi(tokenRegistryAbiStr)
	NameRegistryAbi  = newAbi(nameRegistryAbiStr)
)

func newAbi(abiStr string) *abi.ABI {
	a := &abi.ABI{}
	if err := a.UnmarshalJSON([]byte(abiStr)); err != nil {
		log.Fatal(err)
	}
	return a
}

type TransferEvent struct {
	Sender   common.Address `fieldName:"from" fieldId:"0"`
	Receiver common.Address `fieldName:"to" fieldId:"1"`
	Value    *big.Int       `fieldName:"value" fieldId:"2"`
}

func (e *TransferEvent) ConvertDown() *types.TransferEvent {
	evt := &types.TransferEvent{}
	evt.Sender = e.Sender
	evt.Receiver = e.Receiver
	evt.Value = e.Value

	return evt
}

type ApprovalEvent struct {
	Owner   common.Address `fieldName:"owner" fieldId:"0"`
	Spender common.Address `fieldName:"spender fieldId:"1""`
	Value   *big.Int       `fieldName:"value" fieldId:"2"`
}

func (e *ApprovalEvent) ConvertDown() *types.ApprovalEvent {
	evt := &types.ApprovalEvent{}
	evt.Owner = e.Owner
	evt.Spender = e.Spender
	evt.Value = e.Value

	return evt
}

type RingMinedEvent struct {
	RingIndex     *big.Int       `fieldName:"_ringIndex" fieldId:"0"`
	RingHash      common.Hash    `fieldName:"_ringhash" fieldId:"1" indexed:"true"`
	Miner         common.Address `fieldName:"_miner" fieldId:"2"`
	FeeRecipient  common.Address `fieldName:"_feeRecipient" fieldId:"3"`
	OrderHashList [][32]uint8    `fieldName:"_orderHashList" fieldId:"4"`
	AmountsList   [][6]*big.Int  `fieldName:"_amountsList" fieldId:"5"`
}

func (e *RingMinedEvent) ConvertDown() (*types.RingMinedEvent, []*types.OrderFilledEvent, error) {
	length := len(e.OrderHashList)

	if length != len(e.AmountsList) || length < 2 {
		return nil, nil, errors.New("ringMined event unpack error:orderHashList length invalid")
	}

	evt := &types.RingMinedEvent{}
	evt.RingIndex = e.RingIndex
	evt.Ringhash = e.RingHash
	evt.Miner = e.Miner
	evt.FeeRecipient = e.FeeRecipient

	var list []*types.OrderFilledEvent
	lrcFee := big.NewInt(0)
	for i := 0; i < length; i++ {
		var (
			fill                        types.OrderFilledEvent
			preOrderHash, nextOrderHash common.Hash
		)

		if i == 0 {
			preOrderHash = common.Hash(e.OrderHashList[length-1])
			nextOrderHash = common.Hash(e.OrderHashList[1])
		} else if i == length-1 {
			preOrderHash = common.Hash(e.OrderHashList[length-2])
			nextOrderHash = common.Hash(e.OrderHashList[0])
		} else {
			preOrderHash = common.Hash(e.OrderHashList[i-1])
			nextOrderHash = common.Hash(e.OrderHashList[i+1])
		}

		fill.Ringhash = e.RingHash
		fill.PreOrderHash = preOrderHash
		fill.OrderHash = common.Hash(e.OrderHashList[i])
		fill.NextOrderHash = nextOrderHash

		// [_amountS, _amountB, _lrcReward, _lrcFee, splitS, splitB]. amountS&amountB为单次成交量
		fill.RingIndex = e.RingIndex
		fill.AmountS = e.AmountsList[i][0]
		fill.AmountB = e.AmountsList[i][1]
		fill.LrcReward = e.AmountsList[i][2]
		fill.LrcFee = e.AmountsList[i][3]
		fill.SplitS = e.AmountsList[i][4]
		fill.SplitB = e.AmountsList[i][5]
		fill.FillIndex = big.NewInt(int64(i))

		lrcFee = lrcFee.Add(lrcFee, fill.LrcFee)
		list = append(list, &fill)
	}

	evt.TotalLrcFee = lrcFee
	evt.TradeAmount = length

	return evt, list, nil
}

type OrderCancelledEvent struct {
	OrderHash       common.Hash `fieldName:"_orderHash" fieldId:"0"`
	AmountCancelled *big.Int    `fieldName:"_amountCancelled" fieldId:"1"` // amountCancelled为多次取消累加总量，根据orderhash以及amountCancelled可以确定其唯一性
}

func (e *OrderCancelledEvent) ConvertDown() *types.OrderCancelledEvent {
	evt := &types.OrderCancelledEvent{}
	evt.OrderHash = e.OrderHash
	evt.AmountCancelled = e.AmountCancelled

	return evt
}

type CutoffEvent struct {
	Owner  common.Address `fieldName:"_address" fieldId:"0"`
	Cutoff *big.Int       `fieldName:"_cutoff" fieldId:"1"`
}

func (e *CutoffEvent) ConvertDown() *types.CutoffEvent {
	evt := &types.CutoffEvent{}
	evt.Owner = e.Owner
	evt.Cutoff = e.Cutoff

	return evt
}

type CutoffPairEvent struct {
	Owner  common.Address `fieldName:"_address" fieldId:"0"`
	Token1 common.Address `fieldName:"_token1" fieldId:"1"`
	Token2 common.Address `fieldName:"_token2" fieldId:"2"`
	Cutoff *big.Int       `fieldName:"_cutoff" fieldId:"3"`
}

func (e *CutoffPairEvent) ConvertDown() *types.CutoffPairEvent {
	evt := &types.CutoffPairEvent{}
	evt.Owner = e.Owner
	evt.Token1 = e.Token1
	evt.Token2 = e.Token2
	evt.Cutoff = e.Cutoff

	return evt
}

type TokenRegisteredEvent struct {
	Token  common.Address `fieldName:"addr" fieldId:"0"`
	Symbol string         `fieldName:"symbol" fieldId:"1"`
}

func (e *TokenRegisteredEvent) ConvertDown() *types.TokenRegisterEvent {
	evt := &types.TokenRegisterEvent{}
	evt.Token = e.Token
	evt.Symbol = e.Symbol

	return evt
}

type TokenUnRegisteredEvent struct {
	Token  common.Address `fieldName:"addr" fieldId:"0"`
	Symbol string         `fieldName:"symbol" fieldId:"1"`
}

func (e *TokenUnRegisteredEvent) ConvertDown() *types.TokenUnRegisterEvent {
	evt := &types.TokenUnRegisterEvent{}
	evt.Token = e.Token
	evt.Symbol = e.Symbol

	return evt
}

type AddressAuthorizedEvent struct {
	ContractAddress common.Address `fieldName:"addr" fieldId:"0"`
	Number          int            `fieldName:"number" fieldId:"1"`
}

func (e *AddressAuthorizedEvent) ConvertDown() *types.AddressAuthorizedEvent {
	evt := &types.AddressAuthorizedEvent{}
	evt.Protocol = e.ContractAddress
	evt.Number = e.Number

	return evt
}

type AddressDeAuthorizedEvent struct {
	ContractAddress common.Address `fieldName:"addr" fieldId:"0"`
	Number          int            `fieldName:"number" fieldId:"1"`
}

func (e *AddressDeAuthorizedEvent) ConvertDown() *types.AddressDeAuthorizedEvent {
	evt := &types.AddressDeAuthorizedEvent{}
	evt.Protocol = e.ContractAddress
	evt.Number = e.Number

	return evt
}

// event  Deposit(address indexed dst, uint wad);
type WethDepositEvent struct {
	DstAddress common.Address `fieldName:"dst" fieldId:"0"` // 充值到哪个地址
	Value      *big.Int       `fieldName:"wad" fieldId:"1"`
}

func (e *WethDepositEvent) ConvertDown() *types.WethDepositEvent {
	evt := &types.WethDepositEvent{}
	evt.Value = e.Value

	return evt
}

// event  Withdrawal(address indexed src, uint wad);
type WethWithdrawalEvent struct {
	SrcAddress common.Address `fieldName:"src" fieldId:"0"`
	Value      *big.Int       `fieldName:"wad" fieldId:"1"`
}

func (e *WethWithdrawalEvent) ConvertDown() *types.WethWithdrawalEvent {
	evt := &types.WethWithdrawalEvent{}
	evt.Value = e.Value

	return evt
}

type SubmitRingMethod struct {
	AddressList        [][3]common.Address `fieldName:"addressList" fieldId:"0"`   // owner,tokenS,tokenB(authAddress)
	UintArgsList       [][7]*big.Int       `fieldName:"uintArgsList" fieldId:"1"`  // amountS, amountB, validSince (second),validUntil (second), lrcFee, rateAmountS, and walletId.
	Uint8ArgsList      [][1]uint8          `fieldName:"uint8ArgsList" fieldId:"2"` // marginSplitPercentageList
	BuyNoMoreThanBList []bool              `fieldName:"buyNoMoreThanAmountBList" fieldId:"3"`
	VList              []uint8             `fieldName:"vList" fieldId:"4"`
	RList              [][32]uint8         `fieldName:"rList" fieldId:"5"`
	SList              [][32]uint8         `fieldName:"sList" fieldId:"6"`
	MinerId            *big.Int            `fieldName:"minerId" fieldId:"7"`
	FeeSelections      uint16              `fieldName:"feeSelections" fieldId:"8"`
	Protocol           common.Address
}

// should add protocol, miner, feeRecipient
func (m *SubmitRingMethod) ConvertDown() ([]*types.Order, error) {
	var list []*types.Order
	length := len(m.AddressList)
	vsrLength := 2*length + 1

	if length != len(m.UintArgsList) || length != len(m.Uint8ArgsList) || vsrLength != len(m.VList) || vsrLength != len(m.SList) || vsrLength != len(m.RList) || length < 2 {
		return nil, fmt.Errorf("ringMined method unpack error:orders length invalid")
	}

	for i := 0; i < length; i++ {
		var order types.Order

		order.Protocol = m.Protocol
		order.Owner = m.AddressList[i][0]
		order.TokenS = m.AddressList[i][1]
		if i == length-1 {
			order.TokenB = m.AddressList[0][1]
		} else {
			order.TokenB = m.AddressList[i+1][1]
		}
		order.AuthAddr = m.AddressList[i][2]

		order.AmountS = m.UintArgsList[i][0]
		order.AmountB = m.UintArgsList[i][1]
		order.ValidSince = m.UintArgsList[i][2]
		order.ValidUntil = m.UintArgsList[i][3]
		order.LrcFee = m.UintArgsList[i][4]
		// order.rateAmountS
		order.WalletId = m.UintArgsList[i][6]

		order.MarginSplitPercentage = m.Uint8ArgsList[i][0]

		order.BuyNoMoreThanAmountB = m.BuyNoMoreThanBList[i]

		order.V = m.VList[i]
		order.R = m.RList[i]
		order.S = m.SList[i]

		list = append(list, &order)
	}

	return list, nil
}

type CancelOrderMethod struct {
	AddressList    [4]common.Address `fieldName:"addresses" fieldId:"0"`   //  owner, tokenS, tokenB, authAddr
	OrderValues    [7]*big.Int       `fieldName:"orderValues" fieldId:"1"` //  amountS, amountB, validSince (second), validUntil (second), lrcFee, walletId, and cancelAmount
	BuyNoMoreThanB bool              `fieldName:"buyNoMoreThanAmountB" fieldId:"2"`
	MarginSplit    uint8             `fieldName:"marginSplitPercentage" fieldId:"3"`
	V              uint8             `fieldName:"v" fieldId:"4"`
	R              [32]byte          `fieldName:"r" fieldId:"5"`
	S              [32]byte          `fieldName:"s" fieldId:"6"`
}

// todo(fuk): modify internal cancelOrderMethod and implement related functions
func (m *CancelOrderMethod) ConvertDown() (*types.Order, *big.Int, error) {
	var order types.Order

	order.Owner = m.AddressList[0]
	order.TokenS = m.AddressList[1]
	order.TokenB = m.AddressList[2]

	order.AmountS = m.OrderValues[0]
	order.AmountB = m.OrderValues[1]
	order.ValidSince = m.OrderValues[2]
	order.ValidUntil = m.OrderValues[3]
	order.LrcFee = m.OrderValues[4]
	order.WalletId = m.OrderValues[5]
	cancelAmount := m.OrderValues[6]

	order.BuyNoMoreThanAmountB = bool(m.BuyNoMoreThanB)
	order.MarginSplitPercentage = m.MarginSplit

	order.V = m.V
	order.S = m.S
	order.R = m.R

	return &order, cancelAmount, nil
}

type CutoffMethod struct {
	Cutoff *big.Int `fieldName:"cutoff" fieldId:"0"`
}

func (method *CutoffMethod) ConvertDown() *types.CutoffMethodEvent {
	evt := &types.CutoffMethodEvent{}
	evt.Value = method.Cutoff

	return evt
}

type CutoffPairMethod struct {
	Token1 common.Address `fieldName:"token1" fieldId:"0"`
	Token2 common.Address `fieldName:"token2" fieldId:"1"`
	Cutoff *big.Int       `fieldName:"cutoff" fieldId:"2"`
}

func (method *CutoffPairMethod) ConvertDown() *types.CutoffPairMethodEvent {
	evt := &types.CutoffPairMethodEvent{}
	evt.Value = method.Cutoff
	evt.Token1 = method.Token1
	evt.Token2 = method.Token2

	return evt
}

type WethWithdrawalMethod struct {
	Value *big.Int `fieldName:"wad" fieldId:"0"`
}

func (e *WethWithdrawalMethod) ConvertDown() *types.WethWithdrawalMethodEvent {
	evt := &types.WethWithdrawalMethodEvent{}
	evt.Value = e.Value

	return evt
}

type ApproveMethod struct {
	Spender common.Address `fieldName:"spender" fieldId:"0"`
	Value   *big.Int       `fieldName:"value" fieldId:"1"`
}

func (e *ApproveMethod) ConvertDown() *types.ApproveMethodEvent {
	evt := &types.ApproveMethodEvent{}
	evt.Spender = e.Spender
	evt.Value = e.Value

	return evt
}

// function transfer(address to, uint256 value) public returns (bool);
type TransferMethod struct {
	Receiver common.Address `fieldName:"to" fieldId:"0"`
	Value    *big.Int       `fieldName:"value" fieldId:"1"`
}

func (e *TransferMethod) ConvertDown() *types.TransferMethodEvent {
	evt := &types.TransferMethodEvent{}
	evt.Receiver = e.Receiver
	evt.Value = e.Value

	return evt
}

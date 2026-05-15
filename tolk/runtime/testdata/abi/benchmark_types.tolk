contract RuntimeBenchmarkTypes {
    incomingMessages: ()
    storage: null
    forceAbiExport: (
        BenchmarkSmallInt,
        BenchmarkBigInt,
        BenchmarkSmallUint,
        BenchmarkBigUint,
        BenchmarkVarInt16,
        BenchmarkVarUint32,
        BenchmarkBits24,
        BenchmarkCoins,
        BenchmarkBool,
        BenchmarkCell,
        BenchmarkRemaining,
        BenchmarkAddress,
        BenchmarkOptionalAddress,
        BenchmarkExternalAddress,
        BenchmarkAnyAddress,
        BenchmarkNullableRemaining,
        BenchmarkNullableCell,
        BenchmarkRefInt65,
        BenchmarkEmptyTensor,
        BenchmarkNotEmptyTensor,
        BenchmarkSmallIntKeyMap,
        BenchmarkSmallUintKeyMap,
        BenchmarkBigUintKeyMap,
        BenchmarkBitsKeyMap,
        BenchmarkAddressKeyMap,
        BenchmarkDecUnion
    )
}

type BenchmarkSmallInt = int24
type BenchmarkBigInt = int183
type BenchmarkSmallUint = uint53
type BenchmarkBigUint = uint256
type BenchmarkVarInt16 = varint16
type BenchmarkVarUint32 = varuint32
type BenchmarkBits24 = bits24
type BenchmarkCoins = coins
type BenchmarkBool = bool
type BenchmarkCell = cell
type BenchmarkRemaining = RemainingBitsAndRefs
type BenchmarkAddress = address
type BenchmarkOptionalAddress = address?
type BenchmarkExternalAddress = any_address
type BenchmarkAnyAddress = any_address
type BenchmarkNullableRemaining = RemainingBitsAndRefs?
type BenchmarkNullableCell = cell?
type BenchmarkRefInt65 = Cell<int65>
type BenchmarkEmptyTensor = ()
type BenchmarkNotEmptyTensor = (uint123, bool, coins, (int23, int2?), varint32)
type BenchmarkSmallIntKeyMap = map<int32, bool>
type BenchmarkSmallUintKeyMap = map<uint16, address>
type BenchmarkBigUintKeyMap = map<uint78, cell>
type BenchmarkBitsKeyMap = map<bits16, map<int64, (address, coins)>>
type BenchmarkAddressKeyMap = map<address, coins>
type BenchmarkDecUnion = int16 | int128

fun onInternalMessage() {}

contract RuntimeRefs {
    incomingMessages: ()
    storage: null
    forceAbiExport: (
        AddressAndBalance,
        GoodNaming,
        ManyRefsMsg,
        GoodNamingForMsg,
        Role,
        Operation,
        Operation2,
        Operation3
    )
}

struct AddressAndBalance {
    addr: address
    balance: coins
}

type GoodNaming = AddressAndBalance

enum Role : uint1 {
    Admin = 0,
    Aboba = 1,
}

enum Operation : varuint16 {
    Add = 0,
    Sub = 100000000000000,
}

enum Operation2 : int128 {
    TopUp = -10000,
    Withdraw = 10000,
}

enum Operation3 : uint8 {
    Nothing = 10,
    Something = 1,
}

struct (0xdeadbeef) ManyRefsMsg {
    user1: GoodNaming
    user2: Cell<GoodNaming>?
    user3: cell
    user4: AddressAndBalance? = null
    user5: Cell<AddressAndBalance>
    role: Role
    oper1: Operation
    oper2: Operation2
    oper3: Operation3
}

type GoodNamingForMsg = ManyRefsMsg

fun onInternalMessage() {}

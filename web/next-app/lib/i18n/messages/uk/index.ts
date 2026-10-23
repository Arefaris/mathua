// The Ukrainian dictionary, one JSON file per area of the product.
//
// Split by area so no file grows past what a reviewer can read, and so two people
// translating different screens never edit the same file. The first key of every message is
// the file name: `common.save` lives in `common.json`.
import common from './common.json'
import shell from './shell.json'
import learn from './learn.json'
import goals from './goals.json'
import account from './account.json'
import explore from './explore.json'

export const uk = { common, shell, learn, goals, account, explore }

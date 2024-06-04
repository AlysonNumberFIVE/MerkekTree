# MerkelTree
Also, I realized too late that I spelt it Merkel instead of Merkle :'(


## API reference
This is a repository for a Merket Tree implementation project.  The following operations are implemented:
- Get/Lookup
- Put/Insert
- Update
- GenerateProof
- VerifyProof

### Lookup

```
func (merkelTree *MerkelTree) Lookup(hash []byte) (*Node, error)
```
Returns a pointer to the leaf node with the requested `hash` and `error` if the hash is invalid/data doesn't exist

### Insert
```
func (merkelTree *MerkelTree) Insert(data []byte) ([]byte, error)
```
Creates a new record/node in the tree. It returns a hash if successful and an `error` on failure.<br>
`Insert` must not be used to update existing data. Submitting the same piece of data will throw a `rather use Update()` error
as the hash is created from `data`.<br><br>
In a production Merkel Tree, this duplicate condition wouldn't be possible as data would be tied to a unique entry timestamp, making a duplicate insert impossible (unless it's under malicious intent).

### Update
```
func (merkelTree *MerkelTree) Update(newData, hash []byte) ([]byte, error)
```
`Update` takes in `newData` that will be overwriting the data that exists at `hash` and return a new `hash` for the new data.<br><br>
`Update` has support for stale hashes. If the data at `hash` has been updated more than once, all historical `hash`s that node has always had will be valid for lookup as the lookup structure uses a <a href="https://en.wikibooks.org/wiki/Data_Structures/Hash_Tables">chained hashmap</a> to preserve historical hashes.

### GenerateProof
```
func (merkelTree *MerkelTree) GenerateProof(leafHash []byte) (*MerkelProof, error)
```
`GenerateProof` creates a new proof for the node referenced by `leafHash`. It returns a `MerkelProof` if successful or `error` otherwise; in cases of an invalid hash that fails during inital lookup.

### VerifyProof
```
func VerifyProof(proof *MerkelProof, rootHash []byte) bool
```
Verifies the validity of a `MerkelProof` by using the Merkel Tree verification algorithm.

## Documentation

An accompanying PDF design document accompanies this code and must be treated as a helper for reading this code <b> and not a substitute for going through the code</b>.

## Authored by Alyson Ng

package proofs

import "github.com/celestiaorg/go-square/merkle"

// RowProof is a Merkle proof that a set of rows exist in a Merkle tree with a
// given data root.
type RowProof2 struct {
	// RowRoots are the roots of the rows being proven.
	RowRoots []byte `json:"row_roots"`
	// Proofs is a list of Merkle proofs where each proof proves that a row
	// exists in a Merkle tree with a given data root.
	Proofs   []*merkle.Proof `json:"proofs"`
	StartRow uint32          `json:"start_row"`
	EndRow   uint32          `json:"end_row"`
}

type RowProof struct {
	RowRoots [][]byte `protobuf:"bytes,1,rep,name=row_roots,json=rowRoots,proto3" json:"row_roots,omitempty"`
	Proofs   []*Proof `protobuf:"bytes,2,rep,name=proofs,proto3" json:"proofs,omitempty"`
	Root     []byte   `protobuf:"bytes,3,opt,name=root,proto3" json:"root,omitempty"`
	StartRow uint32   `protobuf:"varint,4,opt,name=start_row,json=startRow,proto3" json:"start_row,omitempty"`
	EndRow   uint32   `protobuf:"varint,5,opt,name=end_row,json=endRow,proto3" json:"end_row,omitempty"`
}

type Proof struct {
	Total    int64    `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Index    int64    `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	LeafHash []byte   `protobuf:"bytes,3,opt,name=leaf_hash,json=leafHash,proto3" json:"leaf_hash,omitempty"`
	Aunts    [][]byte `protobuf:"bytes,4,rep,name=aunts,proto3" json:"aunts,omitempty"`
}

package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	//	"bytes"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	//	"6.5840/labgob"
	"6.5840/labrpc"
)

// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

type Role int

const (
	LEADER Role = iota
	FOLLOWER
	CANDIDATE
)

const (
	HEARTBEAT_INTERVAL = 10 * time.Millisecond
)

// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	role    Role // leader/follower/candidate
	term    int  // current term index
	voteFor int  // vote for peer's index

	electionTimer   *time.Timer // timer of election timeout, random value
	heartbeatsTimer *time.Timer // timer of heartbeat timeout, fixed value (10ms)
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	rf.mu.Lock()
	term = rf.term
	if rf.role == LEADER {
		isleader = true
	} else {
		isleader = false
	}
	rf.mu.Unlock()

	return term, isleader
}

func randomTimeout() time.Duration {
	return time.Duration(150+rand.Int31n(150)) * time.Millisecond
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
// before you've implemented snapshots, you should pass nil as the
// second argument to persister.Save().
// after you've implemented snapshots, pass the current snapshot
// (or nil if there's not yet a snapshot).
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// raftstate := w.Bytes()
	// rf.persister.Save(raftstate, nil)
}

// restore previously persisted state.
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).

}

// example RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term        int
	CandidateId int
}

// example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int
	VoteGranted bool
}

// example RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	if args.Term > rf.term || (args.Term == rf.term && rf.voteFor == -1) {
		rf.role = FOLLOWER
		rf.term = args.Term
		rf.voteFor = args.CandidateId
		rf.electionTimer = time.NewTimer(randomTimeout())
		reply.VoteGranted = true
		reply.Term = rf.term
		rf.mu.Unlock()
		return
	} else {
		reply.VoteGranted = false
		reply.Term = rf.term
		rf.mu.Unlock()
		return
	}

}

// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).

	return index, term, isLeader
}

// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

func (rf *Raft) ticker() {
	for !rf.killed() {

		// Your code here (2A)
		// Check if a leader election should be started.
		select {
		case <-rf.electionTimer.C:
			rf.ElectLeader()
		case <-rf.heartbeatsTimer.C:
			rf.HearBeat()
		}

		// pause for a random amount of time between 50 and 350
		// milliseconds.
		ms := 50 + (rand.Int63() % 300)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

func (rf *Raft) ElectLeader() {

	rf.mu.Lock()
	if rf.role == LEADER {
		rf.mu.Unlock()
		return
	}

	// become candidate
	rf.electionTimer = time.NewTimer(randomTimeout())
	rf.role = CANDIDATE
	rf.term++
	rf.voteFor = rf.me
	voteReceived := 1
	req := RequestVoteArgs{
		CandidateId: rf.me,
		Term:        rf.term,
	}
	rf.mu.Unlock()

	// request all nodes for voting me
	for i := 0; i < len(rf.peers); i++ {

		if i == rf.me {
			continue
		}

		rf.mu.Lock()
		if rf.role != CANDIDATE {
			rf.mu.Unlock()
			break
		}
		rf.mu.Unlock()

		go func(i int) {

			reply := RequestVoteReply{}

			// if this peer doesn't response, continue request next peer
			ok := rf.sendRequestVote(i, &req, &reply)
			if !ok {
				return
			}

			rf.mu.Lock()
			if reply.Term > rf.term { // if find other peer's term greater than me, become follower immediately
				rf.role = FOLLOWER
				rf.term = reply.Term
				rf.voteFor = -1
				rf.mu.Unlock()
				return
			} else if reply.VoteGranted {
				voteReceived++
				// if most of peers vote me, I will be the leader
				if voteReceived > len(rf.peers)/2 && rf.role == CANDIDATE {
					rf.role = LEADER
					rf.voteFor = -1
					rf.mu.Unlock()
					go rf.HearBeat()
					return
				}
				rf.mu.Unlock()
			} else {
				rf.mu.Unlock()
			}

		}(i)

	}

}

type AppendEntriesArgs struct {
	Term int
}

type AppendEntriesReply struct {
	Term int
}

func (rf *Raft) HearBeat() {

	for !rf.killed() {
		// once the peer is dead or it's not leader, the peer will stop sending heartbeat to others
		rf.mu.Lock()
		if rf.role != LEADER {
			rf.mu.Unlock()
			return
		}
		rf.heartbeatsTimer.Reset(HEARTBEAT_INTERVAL)
		args := AppendEntriesArgs{Term: rf.me}

		rf.mu.Unlock()

		wg := sync.WaitGroup{}
		hearBeatNoResp := 0

		// sending heartbeat to other peers
		for i := 0; i < len(rf.peers); i++ {
			if i == rf.me {
				rf.mu.Lock()
				rf.electionTimer = time.NewTimer(randomTimeout())
				rf.mu.Unlock()
				continue
			}

			// once the peer is dead or it's not leader, the peer will stop sending heartbeat to others
			rf.mu.Lock()
			if rf.role != LEADER {
				rf.mu.Unlock()
				break
			}
			rf.mu.Unlock()

			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				reply := AppendEntriesReply{}
				ok := rf.peers[i].Call("Raft.AppendEntries", &args, &reply)
				if !ok {
					rf.mu.Lock()
					hearBeatNoResp++
					rf.mu.Unlock()
					return
				}

				// if find other peer's term greater than me, become follower immediately
				rf.mu.Lock()
				if reply.Term > rf.term {
					rf.term = reply.Term
					rf.role = FOLLOWER
					rf.voteFor = -1
				}
				rf.mu.Unlock()
			}(i)

		}

		wg.Wait()
		rf.mu.Lock()
		if hearBeatNoResp > len(rf.peers)/2 && rf.role == LEADER {
			rf.role = FOLLOWER
			rf.voteFor = -1
			rf.mu.Unlock()
			return
		}
		rf.mu.Unlock()
	}
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {

	rf.mu.Lock()
	reply.Term = rf.term
	// only update my election timer when it's a valid leader
	if args.Term >= rf.term {
		rf.electionTimer = time.NewTimer(randomTimeout())
	}
	rf.mu.Unlock()
}

// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.role = FOLLOWER
	rf.term = 0
	rf.voteFor = -1

	rf.electionTimer = time.NewTimer(randomTimeout())
	rf.heartbeatsTimer = time.NewTimer(HEARTBEAT_INTERVAL)

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.ticker()

	return rf
}

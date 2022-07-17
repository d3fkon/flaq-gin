package users

// Get all waitlisted users. Mostly used by admins
func GetWaitlist() {}

// Register a user to the waitlist
func RegisterToWaitlist() {}

// Get the total waitlist queue count
// Used to show the total number of users in the queue to a user
func GetWaitlistQueue() {}

// Get the waitlist queue cardinality for a user
// Show the cardinality of the user in the queue
func GetWaitlistCardinality() {}

// Allow a user into the application by assigning a referral code to them
func AllowUserByReferralCode() {}

// Allow the user even without a referral code
// Usually used by other functions or by admin
func AllowUserRaw() {}

// Set the base queue count to bolster the queue size
func SetBaseQueueCount() {}

// Email the first N users in the queue
func EmailFirstNUsersInQueue() {}

// Allow the first N users in the queue
// Usually couple by emailing them that they are now allowed
func AllowFirstNUsersInQueue() {}

// Email the last N users in the queue
func EmailLastNUsersInQueue() {}

// Allow the last N users in the queue
// Usuallly coupled with emailing them that they are allowed
func AllowLastNUsersInQueue() {}

// Apply a referral code asssigned to the user
// Check if the referral code provided by the user, matches the referral code assigned to the user
func ApplyReferralCode() {}

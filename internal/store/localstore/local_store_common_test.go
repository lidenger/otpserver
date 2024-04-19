package localstore

import "testing"

func TestManageLocalStoreFile(t *testing.T) {
	err := manageLocalStoreFile("G:\\otpserver\\store", "otp_account_secret")
	if err != nil {
		t.Error(err)
	}
}

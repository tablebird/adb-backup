package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQueryResult(t *testing.T) {
	result := `Row: 0 _id=8, thread_id=1, address=12121212, person=NULL, date=1765200801134, date_sent=1765197954000, protocol=0, read=1, status=-1, type=1, reply_path_present=0, subject=NULL, body=Hi there
Long time no talk
How have you been lately
Hope everything is going well with your new job, service_center=+12121212, locked=0, error_code=0, seen=1, timed=0, deleted=0, sync_state=0, marker=0, source=NULL, bind_id=0, mx_status=0, mx_id=NULL, out_time=0, account=NULL, sim_id=2, block_type=0, advanced_seen=3, b2c_ttl=0, b2c_numbers=NULL, fake_cell_type=0, url_risky_type=0, creator=NULL, favorite_date=0, mx_id_v2=NULL, sub_id=-1
Row: 1 _id=7, thread_id=2, address=51472, person=NULL, date=1765007618426, date_sent=1765007616000, protocol=0, read=1, status=-1, type=1, reply_path_present=0, subject=NULL, body=Hey, are we still meeting at the café at 3 PM today?, service_center=+12121212, locked=0, error_code=0, seen=1, timed=0, deleted=0, sync_state=0, marker=0, source=NULL, bind_id=0, mx_status=0, mx_id=NULL, out_time=0, account=NULL, sim_id=2, block_type=0, advanced_seen=3, b2c_ttl=0, b2c_numbers=NULL, fake_cell_type=0, url_risky_type=0, creator=NULL, favorite_date=0, mx_id_v2=NULL, sub_id=-1
Row: 2 _id=6, thread_id=3, address=ttttt, person=NULL, date=1765007390660, date_sent=1765007388000, protocol=0, read=1, status=-1, type=1, reply_path_present=0, subject=NULL, body=Good morning
The sun is shining today

Perfect weather for a walk in the park
Maybe we can go after lunch, service_center=+12121212, locked=0, error_code=0, seen=1, timed=0, deleted=0, sync_state=0, marker=0, source=NULL, bind_id=0, mx_status=0, mx_id=NULL, out_time=0, account=NULL, sim_id=2, block_type=0, advanced_seen=3, b2c_ttl=0, b2c_numbers=NULL, fake_cell_type=0, url_risky_type=0, creator=NULL, favorite_date=0, mx_id_v2=NULL, sub_id=-1
Row: 3 _id=5, thread_id=1, address=12121212, person=NULL, date=1765007115692, date_sent=1765006071000, protocol=0, read=1, status=-1, type=1, reply_path_present=0, subject=NULL, body=Dear Mr. Smith,
Your bank account has received a deposit of $1,250.50.
Please log in to your online banking portal to verify the transaction.
Thank you for using our services., service_center=+12121212, locked=0, error_code=0, seen=1, timed=0, deleted=0, sync_state=0, marker=0, source=NULL, bind_id=0, mx_status=0, mx_id=NULL, out_time=0, account=NULL, sim_id=2, block_type=0, advanced_seen=3, b2c_ttl=0, b2c_numbers=NULL, fake_cell_type=0, url_risky_type=0, creator=NULL, favorite_date=0, mx_id_v2=NULL, sub_id=-1
Row: 4 _id=1, thread_id=1, address=12121212, person=NULL, date=1765005909456, date_sent=1765005286000, protocol=0, read=1, status=-1, type=1, reply_path_present=0, subject=NULL, body= Reminder for the team,
The project deadline is this Friday, October 18th, at 5:00 PM sharp.
Please submit all final drafts, data sheets, and presentation slides to the project folder by then.
If you have any questions, or need an extension, reach out to the project manager before Wednesday noon., service_center=+12121212, locked=0, error_code=0, seen=1, timed=0, deleted=0, sync_state=0, marker=0, source=NULL, bind_id=0, mx_status=0, mx_id=NULL, out_time=0, account=NULL, sim_id=2, block_type=0, advanced_seen=3, b2c_ttl=0, b2c_numbers=NULL, fake_cell_type=0, url_risky_type=0, creator=NULL, favorite_date=0, mx_id_v2=NULL, sub_id=-1`
	count := 0
	for rowIndex, text := range parseQueryResult(result) {
		assert.Equal(t, count, rowIndex)
		itemMap := parseItem(text)
		assert.Equal(t, 38, len(itemMap))
		count++
	}
	assert.Equal(t, 5, count)
}

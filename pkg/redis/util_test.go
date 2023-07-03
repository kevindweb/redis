package redis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProtocolParse(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		request  string
		expected [][]string
		wantErr  bool
	}{
		{
			name:    "completely invalid",
			request: "hello",
			wantErr: true,
		},
		{
			name:    "invalid error response",
			request: "-noNewLine",
			wantErr: true,
		},
		{
			name:    "invalid argument width",
			request: "*2\r\n*1\r\n$13\r\nnotenough\r\n*1\r\n$2\r\nOK\r\n",
			wantErr: true,
		},
		{
			name:    "invalid batch number",
			request: "*1\r\n*1\r\n$13\r\ntotallyvalids\r\n*1\r\n$2\r\nOK\r\n",
			wantErr: true,
		},
		{
			name:     "single valid delete request no array",
			request:  "$3\r\ndel\r\n$3\r\nkey\r\n",
			expected: [][]string{{"del", "key"}},
		},
		{
			name:     "single valid set request no bulk",
			request:  "1*\r\n$3\r\nSet\r\n$3\r\nkey\r\n$3\r\nval\r\n",
			expected: [][]string{{"set", "key", "val"}},
		},
		{
			name:     "valid simple string",
			request:  "+simplestring\r\n",
			expected: [][]string{{"simplestring"}},
		},
		{
			name:     "valid simple error response",
			request:  "-this is an error\r\n",
			expected: [][]string{{"err", "this is an error"}},
		},
		{
			name:     "valid bulk error no bulk string",
			request:  "*1\r\n*1\r\n-bulkerr\r\n",
			expected: [][]string{{"err", "bulkerr"}},
		},
		{
			name:     "valid ok",
			request:  "*1\r\n*1\r\n$2\r\nOK\r\n",
			expected: [][]string{{"ok"}},
		},
		{
			name:     "many batched requests",
			request:  "*2\r\n*1\r\n$13\r\ntotallyvalids\r\n*1\r\n$2\r\nOK\r\n",
			expected: [][]string{{"totallyvalids"}, {"ok"}},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			str := &tc.request
			requests, err := parseRequests(str)
			require.Equal(t, tc.wantErr, err != nil, "error expected: %t, got err: %v", tc.wantErr, err)
			if tc.wantErr {
				return
			}
			require.Equal(t, tc.expected, requests)
		})
	}
}
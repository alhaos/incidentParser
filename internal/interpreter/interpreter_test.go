package interpreter

import (
	"parser/internal/incidentParser"
	"testing"
)

func Test(t *testing.T) {

	testCases := []struct {
		desc           string
		datafile       string
		expectedResult bool
		expectedError  error
	}{
		{
			desc:           "Total db time",
			datafile:       `C:\tmp\003\data\Incidents\00000000F47C437C89C0A442B838F55DA4B2392904028A01.txt`,
			expectedResult: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			p := incidentParser.NewParser()
			inc, err := p.Parse(tC.datafile)
			if err != nil {
				t.Error(err)
			}
			i := Default()

			result := i.ShouldExclude(
				inc,
			)

			if result != tC.expectedResult {
				t.Errorf("unexpected result for test case [%s] want: [%v], got: [%v]", tC.desc, tC.expectedResult, result)
			}
		})
	}
}

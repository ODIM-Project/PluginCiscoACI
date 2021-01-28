//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

//Package capresponse ...
package capresponse

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	iris "github.com/kataras/iris/v12"
)

// SetErrorResponse will accepts the iris context, error string and status code
// it will set error resopnse to ctx
func SetErrorResponse(ctx iris.Context, statusCode int32, statusMsg, errMsg string, msgArgs []interface{}) {
	resp := common.GeneralError(statusCode, statusMsg, errMsg, msgArgs, nil)
	ctx.StatusCode(int(resp.StatusCode))
	common.SetResponseHeader(ctx, resp.Header)
	ctx.JSON(resp.Body)
}

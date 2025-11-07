# ==================================================================================
#
#       Copyright (c) 2025 Samsung Electronics Co., Ltd. All Rights Reserved.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
# ==================================================================================


# Since, Kind cluster runs in a docker container, therefore it is required to port-forward the svc to connect
kubectl port-forward svc/mme-modelmgmtservice -n traininghost  32007:8082 & PF_PID=$!
sleep 1
python3 -m pytest ../tests/ --maxfail=1 --disable-warnings -q --html=pytest_report.html

# Stop port-forwarding
kill $PF_PID
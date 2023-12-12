.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

.. Copyright (c) 2023 Samsung Electronics Co., Ltd. All Rights Reserved.

User-Guide
==========

.. contents::
   :depth: 3
   :local:


Overview
--------
Model Management Service works with AIML Framework to manage the life cycle of trained AIML models,
such as creating a model, storing the trained model, storing the trained model info.
It exposes REST based API to work with models.

Steps to build and run Model Management Service Standalone
-----------------------------------------------------------

Prerequisites

#. Install go


Steps

.. code:: bash

         git clone "https://gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice.git"
        cd modelmgmtservice

| Update ENV variables in config.env
| Execute below commands
        
.. code:: bash

        export $(< ./config.env)
        go get
        go build -o mme_bin .
        ./mme_bin

Steps to run Model Management Service using AIMLFW deployment scripts
----------------------------------------------------------------------

Follow the steps in this link: `AIMLFW installation guide <https://docs.o-ran-sc.org/projects/o-ran-sc-aiml-fw-aimlfw-dep/en/latest/installation-guide.html>`__

APIs and samples
-----------------

#. Registering a model in Model Management Service
   Sample model-name value is "qos_301"

   .. code:: bash

        curl  -i  -H "Content-Type: application/json"  \
                -X POST \
                -d '{"model-name":"qos_301", "rapp-id": "rapp_1", "meta-info" : {"accuracy":"90","model-type":"timeseries","feature-list":["pdcpBytesDl","pdcpBytesUl"]}}' \
                http://127.0.0.1:32006/registerModel


#. Fetch trained model information from Model Management Service

   .. code:: bash

        curl -X GET  http://127.0.0.1:32006/getModelInfo/qos_301

#. Upload a trained AIML Model to Model Management Service

   .. code:: bash

        curl -F "file=@<MODEL_ZIP_FILE_NAME>" http://127.0.0.1:32006/uploadModel/qos_301

#. Download a trained model from Model Management Service

   .. code:: bash

        curl -X GET http://127.0.0.1:32006/downloadModel/qos_301/model.zip --output model.zip

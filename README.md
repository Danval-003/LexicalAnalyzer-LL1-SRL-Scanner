### LexicalAnalyzer-LL1-SRL-Scanner 📚🔍

📘 **Version 1.0 - API COMPLETE**

The LexicalAnalyzer-LL1-SRL-Scanner API is a robust Go-based application designed to facilitate lexical analysis using Yalex files and an SLR table based on Yalp files. This API empowers text analysis by converting input text into tokens using scanners, which meticulously identify and categorize various elements within the text. These tokens are subsequently utilized by the SLR table to perform syntactical analysis.

#### Routes 🛣️

- **/scan**: Endpoint to analyze text and generate tokens based on predefined scanners. Send a POST request with the text in the request body to utilize this functionality.
  
- **/srl**: Endpoint providing access to the SLR table. Send a GET request to receive the SLR table in JSON format.
  
- **/swagger**: Interactive API documentation available via Swagger. Accessible from your browser to explore endpoints and their respective request/response structures.

#### Frontend Interface (Angular) 🌐

The frontend interface of this application is developed using Angular, offering a user-friendly environment to interact with the LexicalAnalyzer-LL1-SRL-Scanner API. It enables users to input text for scanning and displays the generated tokens, leveraging the robust backend capabilities provided by the Go-based API.

#### License 📜

Distributed under the MIT License. See the LICENSE file for more information.

#### Getting Started 🚀

To run the application:

1. Clone the repository to your local machine.
   
2. Install required dependencies by running `go get` in the project directory.
   
3. Build the application using `go build`.
   
4. Run the application with `./lexical-analyzer`.

Ensure you have the Yalex files for scanners and the Yalp file for the SLR table in their respective directories before starting the application.

This README provides a comprehensive overview of the LexicalAnalyzer-LL1-SRL-Scanner API and its Angular frontend interface. For more detailed information, refer to the API documentation available on Swagger.

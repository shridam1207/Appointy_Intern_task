# Appointy_Intern_task
<h3>
  Front end task is to create an UI web clone of Inshorts App
  </h3><br>
  
Open the main.html file in your browser to view the output.
<h3>
  Backend task is to create an API for the Inshorts App
  </h3><br>
  <h4>
  Procedure to test the API
  </h4><br>
  
 1. Everything is written in the main.go file.
 2. To run the main.go file run command 'go run main.go' in your terminal.
 3. The local host server at 8081 port, change it accordingly.
 4. To make a get request run 'curl -i http://localhost:8081/articles' in terminal.
 5. To make a post request run 'curl -XPOST  http://localhost:8081/articles -d ' { 
    JSON_request_body
} '
6. To make a specific article id GET request run 'curl -i http://localhost:8081/articles/{id} in terminal.
7. To make a word search request run 'curl -i http://localhost:8081/articles/search?q={param1}&..' in terminal
8. I have not added the pagination and unit_testing as i failed to understand it. I may add it later.

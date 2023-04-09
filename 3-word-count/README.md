

## Dependencies
- Echo Framework

### Curl
```curl --location --request GET 'http://localhost:80/beef/summary'```


## WordCount
takes a string as input and returns a WordCount struct containing
the count of each word in the input string. The word count is stored within
the "Beef" field of the WordCount struct. The function removes commas and
periods from the input string before counting the words.
>
Example Response
```  {
      "beef": {
          "Consequat": 1,
          "venison": 1,
          "quis": 3,
          "jowl": 2,
          ...
      }
  }```

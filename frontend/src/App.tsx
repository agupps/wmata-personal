import { useEffect, useState } from 'react'
import './App.css'

type ApiResponse = {
  message: string;
};


function App() {
  const [data, setData] = useState<ApiResponse | null>(null);

  useEffect(() => {
    fetch("http://localhost:8080/busStops")
      .then((res) => res.json())
      .then(setData);
  }, []);

  return (
    <div>
      <h1>Frontend</h1>

      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}

export default App

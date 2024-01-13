import { useState } from 'react';
import reactLogo from './assets/bicycle-svgrepo-com.svg';
import './App.css';
import axios from 'axios';

const App: React.FC = () => {
  const [txtFile, setTxtFile] = useState<File | null>();
  const [errorMsg, setErrorMsg] = useState<string>('');
  const [urls, setURLs] = useState([]);
  // const [errorURLs, setErrorURLs] = useState([]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setTxtFile(e.target.files?.[0] || null);
  };

  const handleUpload = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (txtFile) {
      setErrorMsg('');
      const data = new FormData();
      data.append('file', txtFile);

      axios
        .post('http://localhost:8080/upload_txt', data, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })
        .then((resp) => {
          if (resp.data.urls.length > 0) {
            const urls = resp.data.urls;
            for (const idx in urls) {
              if (urls[idx] !== '') {
                if (parseInt(idx) === 0) {
                  window.open(urls[idx], '_blank');
                } else {
                  setTimeout(() => {
                    window.open(urls[idx], '_blank');
                  }, 5000);
                }
              }
            }
            setURLs(resp.data.urls);
          }
          setTxtFile(null);
          (document.getElementById('fileInput') as HTMLFormElement).value = '';
        })
        .catch((error: Error) => setErrorMsg(error.message));
    }
  };

  return (
    <>
      <div>
        <img src={reactLogo} className="logo react" alt="React logo" />
      </div>
      <h1>Save files from URLs</h1>
      <div className="card">
        <h2>Select a .txt file with URLs</h2>
        <form className="file-form" onSubmit={(e) => handleUpload(e)}>
          <input
            id="fileInput"
            type="file"
            accept=".txt"
            onChange={(e) => handleFileChange(e)}
          />
          <button type="submit" className="save-btn">
            Upload
          </button>
          {errorMsg !== '' && <p className="error-msg">Error: {errorMsg}</p>}
          {urls.length > 0 && (
            <div className="error-msg">
              {urls.map((url, i) => (
                <p key={i}>{url}</p>
              ))}
            </div>
          )}
        </form>
      </div>
    </>
  );
};

export default App;

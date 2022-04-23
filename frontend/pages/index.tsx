import { Button, Container, Heading } from "@chakra-ui/react";
import Layout from "@components/layouts/centered";
import { useEffect, useState } from "react";

class Request<T = any> {
  private url: string;
  private dataState = useState<T | undefined>(undefined);
  private errorState = useState<Error | undefined>(undefined);

  public get data(): T | undefined {
    return this.dataState[0];
  }

  public get error(): Error | undefined {
    return this.errorState[0];
  }

  constructor(url: string) {
    this.url = url;
  }

  private setData(data: T): void {
    this.dataState[1](data);
    this.errorState[1](undefined);
  }

  private setError(error: Error): void {
    this.dataState[1](undefined);
    this.errorState[1](error);
  }
  
  public send(): void {
    fetch(this.url).then(res => {
      if (res.ok) {
        res.json().then(data => this.setData(data));
      } else {
        throw new Error(res.statusText);
      }
    }).catch(err => this.setError(err));
  }
}

export default function Home() {
  // const { data: ping, error } = Request<{ time: string }>("http://localhost:8080/ping");
  const request = new Request<{ message: string, time: string }>("http://localhost:8080/ping");
  useEffect(() => {
    request.send();
  }, []);


  return (
    <Layout meta={undefined}>
   {/* <Container size="lg" padding="0.5rem"> */}
      <Heading size="3xl">Hello Next.js</Heading>
      <Button onClick={() => request.send() } >Hello</Button>

        {request.data && (<>
          <p>Message: {request.data.message}</p>
          <p>Last ping: {(new Date(request.data.time)).toLocaleString()}</p>
        </>)}
        {request.error &&
          <p>Error: {request.error?.message}</p>
        }
    {/* </Container> */}
    </Layout>
  );
}

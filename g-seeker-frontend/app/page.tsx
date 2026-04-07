import Image from "next/image";
import { Button } from "antd";
import ChatWindow from "./pages/ChatWindow";

export default function Home() {
  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <ChatWindow />
    </div>
  );
}

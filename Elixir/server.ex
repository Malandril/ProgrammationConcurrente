#
# Module Logger
#
defmodule Log do
  def message(str) do
    fmt = fn nb -> nb |> Integer.to_string() |> String.pad_leading(2, "0") end

    {{d, m, y}, {h, min, s}} = :calendar.local_time()

    IO.puts("#{fmt.(y)}/#{fmt.(m)}/#{fmt.(d)} " <> "#{fmt.(h)}:#{fmt.(min)}:#{fmt.(s)} #{str}")
  end
end

#
# Module Server
#
defmodule Server do
  def read_line(socket) do
    {:ok, line} = :gen_tcp.recv(socket, 0)
    line
  end

  defp write_line(socket, line) do
    :gen_tcp.send(socket, line)
  end

  def listen_client(socket, sender, name) do
    # Timeout of 10 seconds
    case :gen_tcp.recv(socket, 0, 10000) do
      {:ok, message} ->
        send(sender, {:send_message, "#{name}: #{message}"})
        listen_client(socket, sender, name)

      {:error, :timeout} ->
        send(sender, {:disconnect, socket})
        :gen_tcp.send(socket, "You timed out\n")
        :gen_tcp.close(socket)
        send(sender, {:send_message, "#{name} timed out\n"})

      {:error, _} ->
        send(sender, {:disconnect, socket})
        :gen_tcp.close(socket)
        send(sender, {:send_message, "#{name} disconnected\n"})
    end
  end

  def serve(socket, sender) do
    write_line(socket, "Login: ")
    name = String.replace(read_line(socket), "\n", "")
    send(sender, {:add_client, socket})
    send(sender, {:send_message, "#{name} connected\n"})
    listen_client(socket, sender, name)
  end

  def message_sender(clients) do
    receive do
      {:add_client, client} ->
        message_sender([client | clients])

      {:send_message, message} ->
        Enum.each(clients, fn x -> write_line(x, message) end)
        message_sender(clients)

      {:disconnect, socket} ->
        message_sender(List.delete(clients, socket))

      {:stop, reason} ->
        exit(reason)
    end
  end

  def accept(socket, sender) do
    # Wait a client
    {:ok, client} = :gen_tcp.accept(socket)
    # Serve new client
    spawn(fn -> serve(client, sender) end)
    # Loop back
    accept(socket, sender)
  end

  def run(port) do
    Log.message("Running server on port #{port}")

    {:ok, socket} =
      :gen_tcp.listen(port, [:binary, packet: :line, active: false, reuseaddr: true])

    sender = spawn(fn -> message_sender([]) end)
    accept(socket, sender)
  end

  def main() do
    {options, _, _} = OptionParser.parse(System.argv(), switches: [port: :integer])
    run(if options[:port], do: options[:port], else: 1234)
  end
end

Server.main()

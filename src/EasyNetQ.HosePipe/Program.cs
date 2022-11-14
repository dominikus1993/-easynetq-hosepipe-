using EasyNetQ.HosePipe;
using EasyNetQ.HosePipe.Infrastructure.RabbitMq;

using RabbitMQ.Client;

IHost host = Host.CreateDefaultBuilder(args)
    .ConfigureServices((ctx, services) =>
    {
        services.AddSingleton<IConnection>(Connection.FromUrl(ctx.Configuration["RabbitMq:Url"]));
        services.AddHostedService<Worker>();
    })
    .Build();

await host.RunAsync();

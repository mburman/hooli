//
//  WHOMessageTableViewController.m
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import "AFNetworking/AFNetworking.h"
#import "AFJSONRPCClient/AFJSONRPCClient.h"
#import "WHOMessageTableViewController.h"
#import "WHONewMessageViewController.h"
#import "WHOMessage.h"
#import "WHOMessageCell.h"

@interface WHOMessageTableViewController () <WHOMessageProtocol>

@end

@implementation WHOMessageTableViewController

- (id)initWithStyle:(UITableViewStyle)style WithUserName:(NSString* )userName
{
    self = [super initWithStyle:style];
    if (self) {
        // Custom initialization
        [self.tableView setSeparatorInset:UIEdgeInsetsZero];
        UILabel* titleLabel = [[UILabel alloc] init];
        [titleLabel setText:@"Hooli"];
        [titleLabel setFont:[UIFont fontWithName:@"Superclarendon-BlackItalic" size:25.0]];
        [titleLabel setTextColor:[UIColor colorWithRed:70.0/255 green:235.0/255 blue:150.0/255 alpha:.85]];
//        [titleLabel setAlpha:0.75];
        [titleLabel.layer setShadowColor:[UIColor darkGrayColor].CGColor];
        [titleLabel.layer setShadowOffset:(CGSize) { .width = 1.5, .height = 1.5 }];
        [titleLabel.layer setShadowRadius:1.5];
        [titleLabel.layer setShadowOpacity:.65];
        [titleLabel sizeToFit];
        [self.navigationItem setTitleView:titleLabel];
        self.messages = [NSMutableArray array];
        self.userName = userName;
        
        self.locationManager = [[CLLocationManager alloc] init];
        self.locationManager.desiredAccuracy = kCLLocationAccuracyBest;
        self.locationManager.delegate = self;
        [self.locationManager startUpdatingLocation];
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    
    // Uncomment the following line to preserve selection between presentations.
    // self.clearsSelectionOnViewWillAppear = NO;
    
    // Uncomment the following line to display an Edit button in the navigation bar for this view controller.
    // self.navigationItem.rightBarButtonItem = self.editButtonItem;
    
    // TODO load messages from server here
    
    UIBarButtonItem* newMessageButton = [[UIBarButtonItem alloc] initWithBarButtonSystemItem:UIBarButtonSystemItemAdd target:self action:@selector(newMessage:)];
    self.navigationItem.rightBarButtonItem = newMessageButton;
    
    [self.tableView setRowHeight:150.0];
    [self.tableView registerNib:[UINib nibWithNibName:@"WHOMessageCell" bundle:[NSBundle mainBundle]] forCellReuseIdentifier:@"MessageCell"];
}

- (void)locationManager:(CLLocationManager *)manager didUpdateLocations:(NSArray *)locations {
    CLLocation *currLoc = [locations firstObject];
    self.userLocation = currLoc;
    
}

- (void)newMessage:(id) sender {
    WHONewMessageViewController* form = [[WHONewMessageViewController alloc] init];
    form.delegate = self;
    [self.navigationController pushViewController:form animated:YES];
    
}

- (void)receivedNewMessage:(NSString *)message {
//    NSLog(@"received new message %@",message);
//    WHOMessage* messageObj = [[WHOMessage alloc] initWithMessage:message Author:self.userName Location:self.userLocation];
//    NSString* latitude = [[NSString alloc] initWithFormat:@"%f", self.userLocation.coordinate.latitude];
//    NSString* longitude = [[NSString alloc] initWithFormat:@"%f", self.userLocation.coordinate.longitude];
    //TODO send message to server
    /*
    AFJSONRPCClient* client = [AFJSONRPCClient clientWithEndpointURL:[NSURL URLWithString:@"http://192.168.1.13:9009"]];
    [client invokeMethod:@"ProposerObj.PostMessage" withParameters:@{@"Message" : message, @"Author" : self.userName, @"Latitude" : latitude, @"Longitude" : longitude} success:^(AFHTTPRequestOperation *operation, id responseObject) {
        NSLog(@"succeeded to send request to server");
    } failure:^(AFHTTPRequestOperation *operation, NSError *error) {
        NSLog(@"failed to send request to server with error %@", error);
    }];
     */
    
    //using REST
    NSURL* baseURL = [NSURL URLWithString:@"http://192.168.1.19:9009/proposer/"];
    NSDictionary* parameters = @{@"MessageText" : message, @"Author" : self.userName, @"Latitude" : [NSNumber numberWithDouble: self.userLocation.coordinate.latitude], @"Longitude" : [NSNumber numberWithDouble: self.userLocation.coordinate.longitude]};
    
    AFHTTPSessionManager *manager = [[AFHTTPSessionManager alloc] initWithBaseURL:baseURL];
//    manager.responseSerializer = [AFJSONResponseSerializer serializer];
    manager.responseSerializer = [AFHTTPResponseSerializer serializer];
//    manager.responseSerializer.acceptableContentTypes = [NSSet setWithObject:@"text/plain"];
    manager.requestSerializer = [AFJSONRequestSerializer serializer];
    
    [manager POST:@"messages" parameters:parameters success:^(NSURLSessionDataTask *task, id responseObject) {
        NSLog(@"REST success!");
    } failure:^(NSURLSessionDataTask *task, NSError *error) {
        NSLog(@"REST failure wih error: %@", error);
    }];
}

- (NSString *)distanceBetweenUserAndLocation:(CLLocation *)location {
    CLLocationDistance CLDistance = [location distanceFromLocation:location];
    NSString* distance = [NSString stringWithFormat:@"%.1f miles away",(CLDistance/1609.344)];
    return distance;
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

#pragma mark - Table view data source

- (NSInteger)numberOfSectionsInTableView:(UITableView *)tableView
{
    // Return the number of sections.
    return 1;
}


- (NSInteger)tableView:(UITableView *)tableView numberOfRowsInSection:(NSInteger)section
{
    // Return the number of rows in the section.
    return self.messages.count;
}

- (UITableViewCell *)tableView:(UITableView *)tableView cellForRowAtIndexPath:(NSIndexPath *)indexPath
{
    WHOMessageCell *cell = [tableView dequeueReusableCellWithIdentifier:@"MessageCell" forIndexPath:indexPath];
    WHOMessage* message = [self.messages objectAtIndex:indexPath.row];
    cell.messageLabel.text = message.message;
    cell.authorLabel.text = message.author;
    cell.distanceLabel.text = [self distanceBetweenUserAndLocation:message.location];
    return cell;
}

- (CGFloat)tableView:(UITableView *)tableView heightForRowAtIndexPath:(NSIndexPath *)indexPath {
    return 150.0;
}

/*
// Override to support conditional editing of the table view.
- (BOOL)tableView:(UITableView *)tableView canEditRowAtIndexPath:(NSIndexPath *)indexPath
{
    // Return NO if you do not want the specified item to be editable.
    return YES;
}
*/

/*
// Override to support editing the table view.
- (void)tableView:(UITableView *)tableView commitEditingStyle:(UITableViewCellEditingStyle)editingStyle forRowAtIndexPath:(NSIndexPath *)indexPath
{
    if (editingStyle == UITableViewCellEditingStyleDelete) {
        // Delete the row from the data source
        [tableView deleteRowsAtIndexPaths:@[indexPath] withRowAnimation:UITableViewRowAnimationFade];
    } else if (editingStyle == UITableViewCellEditingStyleInsert) {
        // Create a new instance of the appropriate class, insert it into the array, and add a new row to the table view
    }   
}
*/

/*
// Override to support rearranging the table view.
- (void)tableView:(UITableView *)tableView moveRowAtIndexPath:(NSIndexPath *)fromIndexPath toIndexPath:(NSIndexPath *)toIndexPath
{
}
*/

/*
// Override to support conditional rearranging of the table view.
- (BOOL)tableView:(UITableView *)tableView canMoveRowAtIndexPath:(NSIndexPath *)indexPath
{
    // Return NO if you do not want the item to be re-orderable.
    return YES;
}
*/

/*
#pragma mark - Navigation

// In a storyboard-based application, you will often want to do a little preparation before navigation
- (void)prepareForSegue:(UIStoryboardSegue *)segue sender:(id)sender
{
    // Get the new view controller using [segue destinationViewController].
    // Pass the selected object to the new view controller.
}
*/

@end
